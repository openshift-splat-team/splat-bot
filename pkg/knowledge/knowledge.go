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
	"github.com/slack-go/slack/socketmode"
	"gopkg.in/yaml.v2"
)

const (
	DEFAULT_URL_PROMPT = `This may be a topic that I can help with.

%s`
	DEFAULT_LLM_PROMPT = `Can you provide a short response that attempts to answer this question: `
)

var (
	knowledgeAssets  = []data.KnowledgeAsset{}
	knowledgeEntries = []data.Knowledge{}
)

func IsMatch(asset data.KnowledgeAsset, tokens []string) bool {

	return isTokenMatch(asset.On, util.NormalizeTokens(tokens))
}

func IsStringMatch(asset data.KnowledgeAsset, str string) bool {
	tokens := strings.Split(str, " ")
	return isTokenMatch(asset.On, util.NormalizeTokens(tokens))
}

func isTokenMatch(match data.TokenMatch, tokens map[string]string) bool {
	tokensMatch := true
	or := match.Type == "or"

	if len(match.Tokens) > 0 {
		if or {
			tokensMatch = util.TokensPresentOR(tokens, match.Tokens...)
		} else {
			tokensMatch = util.TokensPresentAND(tokens, match.Tokens...)
		}
	}

	if tokensMatch && len(match.Terms) > 0 {
		satisfied := 0
		for _, term := range match.Terms {
			tokenMatch := isTokenMatch(term, tokens)
			if tokenMatch {
				satisfied++
				if or {
					break
				}
			}
		}
		if or {
			tokensMatch = satisfied > 0
		} else {
			tokensMatch = satisfied == len(match.Terms)
		}
	}

	return tokensMatch
}

func defaultKnowledgeEventHandler(ctx context.Context, client *socketmode.Client, eventsAPIEvent *slackevents.MessageEvent, args []string) ([]slack.MsgOption, error) {
	return defaultKnowledgeHandler(ctx, args)
}

func defaultKnowledgeHandler(ctx context.Context, args []string) ([]slack.MsgOption, error) {
	matches := []data.KnowledgeAsset{}
	normalizedArgs := util.NormalizeTokens(args)

	for _, entry := range knowledgeAssets {
		if isTokenMatch(entry.On, normalizedArgs) {
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
	Callback:       defaultKnowledgeEventHandler,
	DontGlobQuotes: true,
	MessageOfInterest: func(args []string, attribute data.Attributes, channel string) bool {
		for _, enrty := range knowledgeEntries {
			if enrty.MessageOfInterest(args, attribute, channel) {
				return true
			}
		}
		return true
	},
}
