package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

const (
	JIRA_BASE_URL="https://issues.redhat.com"
	BOT_USER_ID="U03JLP91K47"
)

func StringToBlock(message string, useMarkdown bool) []slack.MsgOption {
	return []slack.MsgOption{	
		slack.MsgOptionText(message, useMarkdown),
		slack.MsgOptionEnableLinkUnfurl(),
	}
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
	return strings.Contains(messageText, fmt.Sprintf("<@%s>", BOT_USER_ID))
}

