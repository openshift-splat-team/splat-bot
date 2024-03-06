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
		if attribute.RequireMention {
			markdownBuffer.WriteString("*")
		}
		markdownBuffer.WriteString("\n")
	}
	markdownBuffer.WriteString("* - requires mention of @splat-bot")
	return markdownBuffer.String()
}

var HelpAttributes = Attributes{
	Regex: `\bhelp\b`,
	RequireMention: true,
	Callback: func(eventsAPIEvent *slackevents.MessageEvent, args []string) ([]slack.MsgOption, error) {
		return []slack.MsgOption{
			slack.MsgOptionText(compileHelp(), true),
		}, nil
	},
}
