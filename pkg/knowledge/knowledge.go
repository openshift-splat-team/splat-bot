package knowledge

import (
	"context"
	"fmt"
	"log"
	"strings"
	"unicode"

	"github.com/openshift-splat-team/splat-bot/pkg/commands"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

const (
	DEFAULT_URL_PROMPT = `This may be a topic that I can help with. %s`
	DEFAULT_LLM_PROMPT = `Can you provide a short response that attempts to answer this question: `
)

var (
	knowledgeEntries = []Knowledge{
		MigrationTopicAttributes,
	}
)

// tokensPresent checks if all of the args are present in the argMap
func tokensPresentAND(argMap map[string]string, args ...string) bool {
	matchedArgs := map[string]bool{}
	for _, arg := range args {
		arg = strings.ToLower(arg)
		if _, exists := argMap[arg]; exists {
			matchedArgs[arg] = true
		}
		if len(matchedArgs) == len(args) {
			return true
		}
	}
	return false
}

// tokensPresent checks if any of the args are present in the argMap
func tokensPresentOR(argMap map[string]string, args ...string) bool {
	for _, arg := range args {
		if _, exists := argMap[strings.ToLower(arg)]; exists {
			log.Printf("found token: %s", arg)
			return true
		}
	}
	return false
}

func stripPunctuation(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return r
		}
		return -1
	}, s)
}

// normalizeTokens convert all tokens to lower case for case insensitive matching
func normalizeTokens(args []string) map[string]string {
	normalized := map[string]string{}
	for _, arg := range args {
		if len(arg) == 0 {
			continue
		}
		arg = stripPunctuation(arg)
		normalized[strings.ToLower(arg)] = arg
	}
	return normalized
}

func defaultKnowledgeHandler(ctx context.Context, client *socketmode.Client, eventsAPIEvent *slackevents.MessageEvent, args []string) ([]slack.MsgOption, error) {
	matches := []Knowledge{}
	log.Println("defaultKnowledgeHandler")

	for _, entry := range knowledgeEntries {
		if entry.MessageOfInterest(args, entry.Attributes) {
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
			responseText = fmt.Sprintf("%s\n%s", responseText, strings.Join(match.URLS, "\n"))
		}
		response = append(response, slack.MsgOptionText(responseText, false))
	}
	return response, nil
}

func init() {
	for _, entry := range knowledgeEntries {
		entry.Callback = defaultKnowledgeHandler
		commands.AddCommand(entry.Attributes)
	}
}

// Knowledge defines a peice of knowledge that the bot can respond with
type Knowledge struct {
	commands.Attributes

	// MarkdownPrompt message that is returned when the prompt matches
	MarkdownPrompt string

	// URLS urls to be appended to a response. if MarkdownPrompt isn't defined, URLS will be
	// attached to a reasonable default message.
	URLS []string

	// when true, the message is sent to an LLM to construct an answer.
	InvokeLLM bool
}

var KnowledgeCommandAttributes = commands.Attributes{
	MessageOfInterest: func(args []string, attribute commands.Attributes) bool {
		for _, enrty := range knowledgeEntries {
			if enrty.MessageOfInterest(args, attribute) {
				return true
			}
		}
		return true
	},
}
