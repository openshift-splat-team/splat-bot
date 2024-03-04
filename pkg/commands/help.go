package commands

import (
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

var HelpAttributes = Attributes{
	Regex: `\bhelp\b`,
	Callback: func(eventsAPIEvent slackevents.EventsAPIEvent) ([]slack.MsgOption, error) {
		return []slack.MsgOption{
			slack.MsgOptionText("SPLAT Bot help:\n", false),
			slack.MsgOptionText("jira create:\n", false),
		}, nil
	},
}
