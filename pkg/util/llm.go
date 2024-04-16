package util

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

const (
	PROMPT_RESPONSE_TIMEOUT = time.Second * 120
)

type Prompt string

// GenerateResponse generates a response from an ollama API endpoint
func GenerateResponse(ctx context.Context, prompt string) (string, error) {
	endpoint := os.Getenv("OLLAMA_ENDPOINT")
	if len(endpoint) == 0 {
		return "", errors.New("OLLAMA_ENDPOINT must be exported")
	}

	model := os.Getenv("OLLAMA_MODEL")
	if len(model) == 0 {
		model = "tinyllama"
	}

	llm, err := ollama.New(ollama.WithModel(model), ollama.WithServerURL(endpoint))
	if err != nil {
		log.Fatal(err)
	}

	timedCtx, cancel := context.WithTimeout(ctx, PROMPT_RESPONSE_TIMEOUT)
	defer cancel()
	completion, err := llms.GenerateFromSinglePrompt(timedCtx, llm, prompt)
	if err != nil {
		log.Fatal(err)
	}
	return completion, nil
}

func HandlePrompt(ctx context.Context, prompt Prompt, client SlackClientInterface, evt *slackevents.MessageEvent) (string, error) {
	log.Printf("channel %s/%s\n", evt.Channel, evt.TimeStamp)
	messages, _, _, err := client.GetConversationReplies(&slack.GetConversationRepliesParameters{
		ChannelID: evt.Channel,
		Timestamp: evt.ThreadTimeStamp,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get thread messages: %s", err)
	}

	messages = AnonymizeMessages(messages)

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

	completion, err := GenerateResponse(ctx, buffer.String())
	if err != nil {
		return "", fmt.Errorf("unable to get response from LLM: %v", err)
	}
	return completion, nil
}
