package commands

import (
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func compileHelp() string {
	var markdownBuffer strings.Builder
	markdownBuffer.WriteString("*SPLAT Bot Help*\n")
	markdownBuffer.WriteString("SPLAT Bot provides automation for the team\n")
	for _, attribute := range attributes {
		markdownBuffer.WriteString(attribute.HelpMarkdown)
		markdownBuffer.WriteString("\n")
	}
	return markdownBuffer.String()
}

var HelpAttributes = Attributes{
	Regex: `\bhelp\b`,
	Callback: func(eventsAPIEvent *slackevents.MessageEvent, args []string) ([]slack.MsgOption, error) {
		return []slack.MsgOption{
			slack.MsgOptionText(compileHelp(), true),
		}, nil
	},
}
