package commands

import (
	"context"
	"strings"

	"github.com/openshift-splat-team/splat-bot/data"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func compileHelp() string {
	var markdownBuffer strings.Builder
	markdownBuffer.WriteString("*SPLAT Bot Help*\n")
	markdownBuffer.WriteString("SPLAT Bot provides automation for the team\n")
	for _, attribute := range attributes {
		if attribute.ExcludeFromHelp {
			continue
		}
		markdownBuffer.WriteString(attribute.HelpMarkdown)
		if attribute.RequireMention {
			markdownBuffer.WriteString("*")
		}
		markdownBuffer.WriteString("\n")
	}
	markdownBuffer.WriteString("* - requires mention of @splat-bot")
	return markdownBuffer.String()
}

var HelpAttributes = data.Attributes{
	Commands:       []string{"help"},
	RequireMention: true,
	Callback: func(ctx context.Context, client *socketmode.Client, eventsAPIEvent *slackevents.MessageEvent, args []string) ([]slack.MsgOption, error) {
		return []slack.MsgOption{
			slack.MsgOptionText(compileHelp(), true),
		}, nil
	},
}
