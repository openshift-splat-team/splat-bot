package commands

import (
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

var CreateAttributes = Attributes{
	Regex: `\bjira\s+create\b`,
	Callback: func(eventsAPIEvent slackevents.EventsAPIEvent) ([]slack.MsgOption, error) {
		return nil, nil
	},
}
