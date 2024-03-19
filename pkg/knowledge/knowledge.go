package knowledge

import (
	"context"
	"fmt"
	"strings"

	"github.com/openshift-splat-team/splat-bot/data"
	"github.com/openshift-splat-team/splat-bot/pkg/commands"
	vsphere "github.com/openshift-splat-team/splat-bot/pkg/knowledge/vsphere"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

const (
	DEFAULT_URL_PROMPT = `This may be a topic that I can help with. %s`
	DEFAULT_LLM_PROMPT = `Can you provide a short response that attempts to answer this question: `
)

var (
	knowledgeEntries = []data.Knowledge{
		vsphere.MigrationTopicAttributes,
		vsphere.ODFTopicAttributes,
		vsphere.VSphere67TopicAttributes,
		vsphere.InstallationX509Attributes,
	}
)

func defaultKnowledgeHandler(ctx context.Context, client *socketmode.Client, eventsAPIEvent *slackevents.MessageEvent, args []string) ([]slack.MsgOption, error) {
	matches := []data.Knowledge{}

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
		entry.DontGlobQuotes = true
		commands.AddCommand(entry.Attributes)
	}
}

var KnowledgeCommandAttributes = data.Attributes{
	MessageOfInterest: func(args []string, attribute data.Attributes) bool {
		for _, enrty := range knowledgeEntries {
			if enrty.MessageOfInterest(args, attribute) {
				return true
			}
		}
		return true
	},
}
