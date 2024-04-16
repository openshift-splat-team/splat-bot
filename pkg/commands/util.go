package commands

import (
	"fmt"

	"github.com/slack-go/slack"
)

const (
	JIRA_BASE_URL = "https://issues.redhat.com"
)

func StringsToBlockUnfurl(messages []string, useMarkdown, unfurlLinks bool) []slack.MsgOption {
	messageBlocks := []slack.Block{}

	for _, message := range messages {
		messageBlocks = append(messageBlocks,
			slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", message, false, false),
				nil,
				nil,
			),
			slack.NewDividerBlock(),
		)

	}

	/*if unfurlLinks {
		messageBlocks = append(messageBlocks, slack.MsgOptionEnableLinkUnfurl())
	}*/
	return []slack.MsgOption{
		slack.MsgOptionBlocks(messageBlocks...),
	}
}

func StringsToBlockWithURLs(messages []string, urls []string) []slack.MsgOption {
	messageBlocks := []slack.Block{}

	for _, message := range messages {
		messageBlocks = append(messageBlocks,
			slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", message, false, false),
				nil,
				nil,
			),
		)
	}

	if len(urls) > 0 {
		messageBlocks = append(messageBlocks, slack.NewDividerBlock())
		messageBlocks = append(messageBlocks, slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "*Relevant links:*", false, false), nil, nil))
	}
	for _, url := range urls {
		messageBlocks = append(messageBlocks,
			slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", url, false, false),
				nil,
				nil,
			),
		)
	}
	return []slack.MsgOption{
		slack.MsgOptionBlocks(messageBlocks...),
	}
}

func StringToBlockUnfurl(message string, useMarkdown, unfurlLinks bool) []slack.MsgOption {
	options := []slack.MsgOption{
		slack.MsgOptionText(message, useMarkdown),
	}
	if unfurlLinks {
		options = append(options, slack.MsgOptionEnableLinkUnfurl())
	}
	return options
}

func StringToBlock(message string, useMarkdown bool) []slack.MsgOption {
	return StringToBlockUnfurl(message, useMarkdown, true)
}

func WrapErrorToBlock(err error, message string) []slack.MsgOption {
	return StringToBlock(fmt.Sprintf("%s: %v", message, err), false)
}
