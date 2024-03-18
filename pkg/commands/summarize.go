package commands

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/openshift-splat-team/splat-bot/pkg/util"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type Prompt string

var (
	PROMPT_ISSUE_TITLE   = Prompt("can you summarize this thread to a single line? The line should be less than 100 characters. ")
	PROMPT_ISSUE_SUMMARY = Prompt("can you summarize this thread a short paragraph?")
)

var SummarizeAttributes = Attributes{
	Regex:          `summary`,
	RequireMention: true,
	RespondInDM:    false,
	Callback: func(ctx context.Context, client *socketmode.Client, evt *slackevents.MessageEvent, args []string) ([]slack.MsgOption, error) {
		response, err := handlePrompt(ctx, PROMPT_ISSUE_SUMMARY, client, evt)
		if err != nil {
			return nil, fmt.Errorf("unable to get summary: %v", err)
		}
		return StringToBlock(fmt.Sprintf("WIP: will summarize thread: %s", response), false), nil
	},
	RequiredArgs: 1,
	HelpMarkdown: "summarize this thread: `summary`",
}

func handlePrompt(ctx context.Context, prompt Prompt, client *socketmode.Client, evt *slackevents.MessageEvent) (string, error) {
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
