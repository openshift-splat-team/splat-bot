package commands

import (
	"context"
	"fmt"

	"github.com/openshift-splat-team/jira-bot/cmd/issue"
	"github.com/openshift-splat-team/splat-bot/data"
	"github.com/openshift-splat-team/splat-bot/pkg/util"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

const issueTemplate = `
*User Story:*
As a ___ I want ___ so I can ___.

*Description:*
< Record any background information >

*Acceptance Criteria:*
< Record how we'll know we're done >

*Other Information:*
< Record anything else that may be helpful to someone else picking up the card >

issue created by splat-bot
`

var CreateAttributes = data.Attributes{
	Commands:       []string{"jira", "create"},
	RequireMention: true,
	Callback: func(ctx context.Context, client util.SlackClientInterface, evt *slackevents.MessageEvent, args []string) ([]slack.MsgOption, error) {
		url := util.GetThreadUrl(evt)
		fmt.Printf("%v", args)
		summary := args[2]
		description := issueTemplate

		if len(url) > 0 {
			description = fmt.Sprintf("%s\n\ncreated from thread: %s", issueTemplate, url)
		}

		issue, err := issue.CreateIssue("SPLAT", summary, description, "Task")
		if err != nil {
			return util.WrapErrorToBlock(err, "error creating issue"), nil
		}
		issueKey := issue.Key
		issueURL := fmt.Sprintf("%s/browse/%s", JIRA_BASE_URL, issueKey)
		return util.StringToBlock(fmt.Sprintf("issue <%s|%s> created", issueURL, issueKey), false), nil
	},
	RequiredArgs: 3,
	HelpMarkdown: "create a Jira issue: `jira create \"[description]\"`",
	ShouldMatch: []string{
		"jira create description",
		"jira create description",
	},
	ShouldntMatch: []string{
		"jira create-with-summary PROJECT bug",
		"jira create-with-summary PROJECT Todo",
	},
}
