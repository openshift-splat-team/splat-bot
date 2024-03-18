package commands

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

const (
	JIRA_BASE_URL = "https://issues.redhat.com"
)

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

func GetThreadUrl(event *slackevents.MessageEvent) string {
	if event.ThreadTimeStamp != "" {
		workspace := "redhat-internal" // Replace with your Slack workspace name
		threadURL := fmt.Sprintf("https://%s.slack.com/archives/%s/p%s",
			workspace, event.Channel, strings.Replace(event.ThreadTimeStamp, ".", "", 1))

		return threadURL
	}
	log.Println("This is not a threaded message")
	return ""
}

func ContainsBotMention(messageText string) bool {
	userID, ok := os.LookupEnv("SPLAT_BOT_USER_ID")
	if !ok {
		log.Println("no bot user id specified with SPLAT_BOT_USER_ID")
		return false
	}
	return strings.Contains(messageText, fmt.Sprintf("<@%s>", userID))
}
