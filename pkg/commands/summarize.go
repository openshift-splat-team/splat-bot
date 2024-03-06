package commands

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

var SummarizeAttributes = Attributes{
	Regex: `\bsummary\b`,
	RequireMention: true,
	RespondInDM: true,
	Callback: func(client *socketmode.Client, evt *slackevents.MessageEvent, args []string) ([]slack.MsgOption, error) {
		response, err := getSummary(client, evt)
		if err != nil {
			return nil, fmt.Errorf("unable to get summary: %v", err)
		}
		return StringToBlock(fmt.Sprintf("WIP: will summarize thread: %s", response), false), nil
	},
	RequiredArgs: 1,
	HelpMarkdown: "summarize this thread: `summary`",
}

func getSummary(client *socketmode.Client, evt *slackevents.MessageEvent) (string, error){
	endpoint := os.Getenv("OLLAMA_ENDPOINT")
	if len(endpoint) == 0 {
		return "", errors.New("OLLAMA_ENDPOINT must be exported")
	}
	llm, err := ollama.New(ollama.WithModel("llama2"), ollama.WithServerURL(endpoint))
	if err != nil {
		log.Fatal(err)
	}

	messages, _, _, err := client.GetConversationReplies(&slack.GetConversationRepliesParameters{
		ChannelID: evt.Channel,
		Timestamp: evt.ThreadTimeStamp,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get thread messages: %s", err)
	}

	buffer := strings.Builder{}
	buffer.WriteString("summarize this thread:\n")
	for _, message := range messages {
		text := message.Msg.Text
		if ContainsBotMention(text) {
			continue
		}
		buffer.WriteString(message.Msg.Text)
		buffer.WriteString("\n")
	}
	log.Printf("sending context to ollama: %s", buffer.String())
	ctx := context.Background()
	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, buffer.String())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Response:\n", completion)
	return completion, nil
}