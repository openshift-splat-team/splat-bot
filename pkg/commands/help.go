package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/openshift-splat-team/splat-bot/data"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func compileHelp() slack.MsgOption {
	messageBlocks := []slack.Block{}
	messageBlocks = append(messageBlocks, slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", "*Command* | *Description*", false, false),
		nil,
		nil,
	))
	for _, attribute := range getAttributes() {
		if attribute.ExcludeFromHelp {
			continue
		}
		messageBlocks = append(messageBlocks, slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("%s | %s", strings.Join(attribute.Commands, " "), attribute.HelpMarkdown), false, false),
			nil,
			nil,
		))

	}
	return slack.MsgOptionBlocks(messageBlocks...)
}

var HelpAttributes = data.Attributes{
	Commands:       []string{"help"},
	RequireMention: true,
	Callback: func(ctx context.Context, client *socketmode.Client, eventsAPIEvent *slackevents.MessageEvent, args []string) ([]slack.MsgOption, error) {
		return []slack.MsgOption{
			compileHelp(),
		}, nil
	},
}
