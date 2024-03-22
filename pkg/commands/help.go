package commands

import (
	"context"
	"strings"

	"github.com/openshift-splat-team/splat-bot/data"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func compileHelp() slack.MsgOption {
	messageBlocks := []slack.Block{}
	for _, attribute := range getAttributes() {
		if attribute.ExcludeFromHelp {
			continue
		}
		messageBlocks = append(messageBlocks,
			slack.NewSectionBlock(
				slack.NewTextBlockObject("plain_text", strings.Join(attribute.Commands, " "), false, false),
				[]*slack.TextBlockObject{
					//slack.NewTextBlockObject("mrkdwn", strings.Join(attribute.Commands, " "), false, false),
					slack.NewTextBlockObject("plain_text", attribute.HelpMarkdown, false, false),
				},
				nil,
			),
			slack.NewDividerBlock(),
		)
	}
	return slack.MsgOptionBlocks(messageBlocks...)
}

var HelpAttributes = data.Attributes{
	Commands:        []string{"help"},
	RequireMention:  true,
	ExcludeFromHelp: true,
	Callback: func(ctx context.Context, client *socketmode.Client, eventsAPIEvent *slackevents.MessageEvent, args []string) ([]slack.MsgOption, error) {
		return []slack.MsgOption{
			compileHelp(),
		}, nil
	},
	ResponseIsEphemeral: true,
	RespondInChannel:    true,
	ShouldMatch: []string{
		"help",
	},
	ShouldntMatch: []string{
		"jira create-with-summary PROJECT bug",
		"jira create-with-summary PROJECT Todo",
	},
}
