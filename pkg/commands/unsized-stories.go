package commands

import (
	"fmt"
	"strings"

	"github.com/openshift-splat-team/jira-bot/pkg/util"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

var UnsizedAttributes = Attributes{
	Regex: `\bjira\s+unsized\b`,
	RequireMention: true,
	Callback: func(evt *slackevents.MessageEvent, args []string) ([]slack.MsgOption, error) {
		issues, err := util.GetUnsizedStories()
		if err != nil {
			return WrapErrorToBlock(err, "error querying issues"), nil
		}

		var builder strings.Builder
		for _, issue := range issues {
			builder.WriteString(fmt.Sprintf("%s - %s\n", issue.Key, issue.Fields.Summary))
		}
		if len(issues) == 0 {
			builder.WriteString("no issues found")
		}

		return StringToBlock(builder.String(), false), nil
	},
	RequiredArgs: 3,
	HelpMarkdown: "outputs a list of unsized stories for import in to PlanIt Poker: `jira unsized [project]`",
}


// LLM hey SPLAT bot, what is going on with Azure
