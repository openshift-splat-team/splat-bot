package commands

import (
	"context"
	"strings"

	"github.com/openshift-splat-team/splat-bot/data"
	"github.com/openshift-splat-team/splat-bot/pkg/util"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func compileHelp() slack.MsgOption {
	helpText := strings.Builder{}
	for _, attribute := range getAttributes() {
		if attribute.ExcludeFromHelp {
			continue
		}
		helpText.WriteString("- ")
		helpText.WriteString(attribute.HelpMarkdown)
		helpText.WriteString("\n")
	}

	return util.StringsToBlockUnfurl([]string{helpText.String()}, false, false)[0]
}

var HelpAttributes = data.Attributes{
	Commands:        []string{"help"},
	RequireMention:  true,
	ExcludeFromHelp: true,
	Callback: func(ctx context.Context, client util.SlackClientInterface, eventsAPIEvent *slackevents.MessageEvent, args []string) ([]slack.MsgOption, error) {
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
