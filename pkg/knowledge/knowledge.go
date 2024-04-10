package knowledge

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/openshift-splat-team/splat-bot/data"
	"github.com/openshift-splat-team/splat-bot/pkg/commands"
	"github.com/openshift-splat-team/splat-bot/pkg/knowledge/platforms"
	"github.com/openshift-splat-team/splat-bot/pkg/util"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"go.uber.org/thriftrw/ptr"
	"gopkg.in/yaml.v2"
)

const (
	DEFAULT_URL_PROMPT = `This may be a topic that I can help with.

%s`
	DEFAULT_LLM_PROMPT       = `Can you provide a short response that attempts to answer this question: `
	DEBUG_CONDITION_MATCHING = false
)

var (
	knowledgeAssets  = []data.KnowledgeAsset{}
	knowledgeEntries = []data.Knowledge{}
	channelIDMap     = map[string]string{}
	slackClient      util.SlackClientInterface
)

func DumpMatchTree(match data.TokenMatch, depth *int64, messages []string) []string {
	if depth == nil {
		depth = ptr.Int64(0)
	} else {
		*depth++
	}
	if messages == nil {
		messages = []string{}
	}
	padding := strings.Repeat(" ", int(*depth))
	matchType := "OR"
	if len(match.Type) > 0 {
		matchType = match.Type
	}
	messages = append(messages, fmt.Sprintf("%s Match: %t; Match Type: %s", padding, match.Satisfied, matchType))
	messages = append(messages, fmt.Sprintf("%s Immediate Tokens: %s", padding, strings.Join(match.Tokens, ",")))
	if len(match.Terms) > 0 {
		messages = append(messages, fmt.Sprintf("%s Number of Descendant Terms(all terms must match in addition to tokens): %d", padding, len(match.Terms)))
		for _, term := range match.Terms {
			messages = DumpMatchTree(term, depth, messages)
		}
	}
	*depth--
	return messages
}

func getCachedClient() (util.SlackClientInterface, error) {
	if slackClient == nil {
		return util.GetClient()
	}
	return slackClient, nil
}

func IsMatch(asset data.KnowledgeAsset, tokens []string) bool {
	if DEBUG_CONDITION_MATCHING {
		log.Printf("+++++++++++++++++++++++++++++++++++++++IsMatch")
		log.Printf("checking if message: %s is relevant to %s", strings.Join(tokens, " "), asset.Name)

		defer func() {
			log.Printf("---------------------------------------IsMatch")
		}()
	}
	return isTokenMatch(&asset.On, util.NormalizeTokens(tokens))
}

func IsStringMatch(asset data.KnowledgeAsset, str string) bool {
	if DEBUG_CONDITION_MATCHING {
		log.Printf("+++++++++++++++++++++++++++++++++++++++IsStringMatch")
		defer func() {
			log.Printf("---------------------------------------IsStringMatch")
		}()

		log.Printf("checking if message: %s is relevant to %s", str, asset.Name)
	}
	tokens := strings.Split(str, " ")
	return IsMatch(asset, tokens)
}

var depth = 0

func isTokenMatch(match *data.TokenMatch, tokens map[string]string) bool {
	var padding string
	if DEBUG_CONDITION_MATCHING {
		depth++
		padding := strings.Repeat("  ", depth)
		log.Printf("%s+isTokenMatch", padding)
	}
	tokensMatch := true
	or := match.Type == "or"

	if len(match.Tokens) > 0 {
		if or {
			tokensMatch = util.TokensPresentOR(tokens, match.Tokens...)
			if DEBUG_CONDITION_MATCHING {
				log.Printf("%sdo any tokens match? %v", padding, tokensMatch)
			}
		} else {
			tokensMatch = util.TokensPresentAND(tokens, match.Tokens...)
			if DEBUG_CONDITION_MATCHING {
				log.Printf("%sdo all tokens match? %v", padding, tokensMatch)
			}
		}
	}

	if DEBUG_CONDITION_MATCHING {
		log.Printf("%stokensMatch: %t; number of match terms: %d", padding, tokensMatch, len(match.Terms))
	}
	if tokensMatch && len(match.Terms) > 0 {
		satisfied := 0
		for idx := range match.Terms {
			tokenMatch := isTokenMatch(&match.Terms[idx], tokens)
			if tokenMatch {
				satisfied++
				if DEBUG_CONDITION_MATCHING {
					log.Printf("%s+term satisfied: %d", padding, satisfied)
				}
				if or {
					if DEBUG_CONDITION_MATCHING {
						log.Printf("%s+or term satisfied", padding)
					}
					satisfied = len(match.Terms)
					break
				}
			}
		}
		tokensMatch = satisfied == len(match.Terms)
		if DEBUG_CONDITION_MATCHING {
			log.Printf("%sall terms satisfied? %v", padding, tokensMatch)
		}
	}
	if DEBUG_CONDITION_MATCHING {
		log.Printf("%s-tokensMatch: %t", padding, tokensMatch)
		depth--
	}
	match.Satisfied = tokensMatch
	return tokensMatch
}

func defaultKnowledgeEventHandler(ctx context.Context, client util.SlackClientInterface, eventsAPIEvent *slackevents.MessageEvent, args []string) ([]slack.MsgOption, error) {
	return defaultKnowledgeHandler(ctx, args, eventsAPIEvent)
}

func getChannelName(channelID string) (string, error) {
	slackClient, err := getCachedClient()
	if err != nil {
		return "", fmt.Errorf("unable to get client: %v", err)
	}
	if name, ok := channelIDMap[channelID]; ok {
		return name, nil
	}

	channel, err := slackClient.GetConversationInfo(
		&slack.GetConversationInfoInput{
			ChannelID: channelID,
		},
	)
	if err != nil {
		return "", fmt.Errorf("error getting channel info: %v", err)
	}
	channelIDMap[channelID] = channel.Name
	return channel.Name, nil
}

func defaultKnowledgeHandler(ctx context.Context, args []string, eventsAPIEvent *slackevents.MessageEvent) ([]slack.MsgOption, error) {
	var channel string
	var err error
	matches := []data.KnowledgeAsset{}

	for idx, entry := range knowledgeAssets {
		if !entry.WatchThreads && eventsAPIEvent.ThreadTimeStamp != "" {
			continue
		}
		if entry.ChannelContext != nil {
			if channel == "" {
				channel, err = getChannelName(eventsAPIEvent.Channel)
				if err != nil {
					return nil, fmt.Errorf("error getting channel name: %v", err)
				}
			}
			channelContext := entry.ChannelContext
			for _, allowedChannel := range channelContext.Channels {
				if channel == allowedChannel {
					terms := platforms.GetPathContextTerms(channelContext.ContextPath)
					for _, term := range terms {
						args = append(args, term.Tokens...)
					}
					break
				}
			}
		}
		if len(entry.RequireInChannel) > 0 {
			if channel == "" {
				channel, err = getChannelName(eventsAPIEvent.Channel)
				if err != nil {
					return nil, fmt.Errorf("error getting channel name: %v", err)
				}
			}
			allowed := false
			for _, requiredChannel := range entry.RequireInChannel {
				if channel == requiredChannel {
					allowed = true
					break
				}
			}
			if !allowed {
				continue
			}
		}
		if isTokenMatch(&knowledgeAssets[idx].On, util.NormalizeTokens(args)) {
			matches = append(matches, entry)
		}
	}

	response := []slack.MsgOption{}
	// TO-DO: how can we handle multiple matches? for now we'll just use the first one
	if len(matches) > 0 {
		match := matches[0]
		// TO-DO: add support for LLM invocation
		//if match.InvokeLLM {}

		responseText := fmt.Sprintf(DEFAULT_URL_PROMPT, match.MarkdownPrompt)

		if len(match.URLS) > 0 {
			//response = append(response, slack.MsgOptionText(strings.Join(match.URLS, "\n"), false))
			response = append(response, commands.StringsToBlockWithURLs([]string{responseText}, match.URLS)...)
		} else {
			response = append(response, slack.MsgOptionText(responseText, true))
		}

	}
	return response, nil
}

func getKnowledgeEntryPaths(path string, paths []string) ([]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("error reading knowledge prompts directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			paths, err = getKnowledgeEntryPaths(filepath.Join(path, file.Name()), paths)
			if err != nil {
				return nil, err
			}
		} else if filepath.Ext(file.Name()) == ".yaml" {
			paths = append(paths, filepath.Join(path, file.Name()))
			continue
		}
	}
	return paths, nil
}

func loadKnowledgeEntries(dir string) error {
	files, err := getKnowledgeEntryPaths(dir, []string{})
	if err != nil {
		return fmt.Errorf("error reading knowledge prompts directory: %v", err)
	}

	for _, filePath := range files {
		log.Printf("loading knowledge entry from %s", filePath)
		knowledgeModel, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("error reading file %s: %v", filePath, err)
		}
		var asset data.KnowledgeAsset
		err = yaml.Unmarshal([]byte(knowledgeModel), &asset)
		if err != nil {
			log.Printf("error unmarshalling file %s: %v", filePath, err)
			continue
		}
		// if the name of a known platform appears in the path add platform specific terms
		// to 'On' which must be met before the knowledge asset is considered a match
		if contextTerms := platforms.GetPathContextTerms(filePath); contextTerms != nil {
			asset.On.Terms = append(asset.On.Terms, contextTerms...)
		}
		knowledgeAssets = append(knowledgeAssets, asset)
	}

	return nil
}

func init() {
	promptPath := os.Getenv("PROMPT_PATH")
	if promptPath == "" {
		promptPath = "/usr/src/app/knowledge_prompts"
	}
	err := loadKnowledgeEntries(promptPath)
	// TODO: Need way for local developers to be able to still start application if they are not testing knowledge stuff.
	//       For now, we will disable the commands tha require this.
	if err != nil {
		fmt.Printf("error loading knowledge entries: %v", err)
		fmt.Println("Skipping adding of knowledge-based actions.")
		return
	}
	commands.AddCommand(KnowledgeCommandAttributes)
}

var KnowledgeCommandAttributes = data.Attributes{
	Callback:           defaultKnowledgeEventHandler,
	DontGlobQuotes:     true,
	RequireMention:     false,
	AllowNonSplatUsers: true,
	MessageOfInterest: func(args []string, attribute data.Attributes, channel string) bool {
		for _, enrty := range knowledgeEntries {
			if enrty.MessageOfInterest(args, attribute, channel) {
				return true
			}
		}
		return true
	},
}
