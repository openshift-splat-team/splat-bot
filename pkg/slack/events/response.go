package events

import (
	"strings"

	"github.com/openshift/ci-tools/pkg/slack/modals"
	"github.com/slack-go/slack"
)

type ResponseOperator int

const (
	RESPONSE_OPERATOR_OR ResponseOperator = iota
	RESPONSE_OPERATOR_AND
)

type AutoResponseStruct struct {
	Keywords           []string
	Response           string
	RequiresChannelTag bool
	Operator           ResponseOperator
}

func ResponseFor(message string, channelMatch bool, responses []AutoResponseStruct) []slack.Block {
	type interaction struct {
		identifier              modals.Identifier
		description, buttonText string
	}

	var blocks []slack.Block

	responseMessage := ""
	lowerMessage := strings.ToLower(message)
	for _, response := range responses {
		if response.RequiresChannelTag && channelMatch == false {
			continue
		}
		matches := 0
		for _, keyword := range response.Keywords {
			if strings.Contains(lowerMessage, strings.ToLower(keyword)) {
				matches = matches + 1
			}
		}
		if response.Operator == RESPONSE_OPERATOR_OR {
			if matches != 0 {
				responseMessage = response.Response
			}
		} else {
			if response.Operator == RESPONSE_OPERATOR_AND {
				if matches == len(response.Keywords) {
					responseMessage = response.Response
				}
			}
		}
	}

	if responseMessage == "" {
		return blocks
	}

	blocks = append(blocks, &slack.SectionBlock{
		Type: slack.MBTSection,
		Text: &slack.TextBlockObject{
			Type: slack.MarkdownType,
			//Text: "Sorry, I don't know how to help with that. Here are all the things I know how to do:",
			Text: responseMessage,
		},
	})
	return blocks
}
