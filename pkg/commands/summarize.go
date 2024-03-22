package commands

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/openshift-splat-team/splat-bot/data"
	"github.com/openshift-splat-team/splat-bot/pkg/util"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type Prompt string

var (
	PROMPT_ISSUE_TITLE   = Prompt("can you summarize this thread to a single line? The line should be less than 100 characters. ")
	PROMPT_ISSUE_SUMMARY = Prompt("provide a brief summary of the thread. only reply with information from the thread: ")
)

var SummarizeAttributes = data.Attributes{
	Commands:       []string{"summary"},
	RequireMention: true,
	RespondInDM:    false,
	Callback: func(ctx context.Context, client util.SlackClientInterface, evt *slackevents.MessageEvent, args []string) ([]slack.MsgOption, error) {
		response, err := handlePrompt(ctx, PROMPT_ISSUE_SUMMARY, client, evt)
		if err != nil {
			return nil, fmt.Errorf("unable to get summary: %v", err)
		}
		return StringToBlock(fmt.Sprintf("Sure! Here is a summary of this thread.\n\n*Note: I am a bot and I try my best to provide a reasonable summary. Be sure to check the summary for accuracy.*\n\n%s\n", response), false), nil
	},
	RequiredArgs: 1,
	HelpMarkdown: "summarize this thread: `summary`",
	ShouldMatch: []string{
		"summary",
	},
	ShouldntMatch: []string{
		"jira create-with-summary PROJECT bug",
		"jira create-with-summary PROJECT Todo",
	},
}

func handlePrompt(ctx context.Context, prompt Prompt, client util.SlackClientInterface, evt *slackevents.MessageEvent) (string, error) {
	log.Printf("channel %s/%s\n", evt.Channel, evt.TimeStamp)
	messages, _, _, err := client.GetConversationReplies(&slack.GetConversationRepliesParameters{
		ChannelID: evt.Channel,
		Timestamp: evt.ThreadTimeStamp,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get thread messages: %s", err)
	}

	buffer := strings.Builder{}
	buffer.WriteString(string(prompt))
	buffer.WriteString("\n")
	for _, message := range messages {
		text := message.Msg.Text
		if ContainsBotMention(text) {
			continue
		}
		buffer.WriteString(message.Msg.Text)
		buffer.WriteString("\n")
	}

	completion, err := util.GenerateResponse(ctx, buffer.String())
	if err != nil {
		return "", fmt.Errorf("unable to get response from LLM: %v", err)
	}
	return completion, nil
}
