package commands

import (
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

var SummarizeAttributes = Attributes{
	Regex: `\bsummary\b`,
	RequireMention: true,
	Callback: func(evt *slackevents.MessageEvent, args []string) ([]slack.MsgOption, error) {
		return StringToBlock("WIP: will summarize thread", false), nil
	},
	RequiredArgs: 2,
	HelpMarkdown: "summary this thread: `summary`",
}
