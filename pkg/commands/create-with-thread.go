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

var CreateJiraWithThreadAttributes = data.Attributes{
	Commands:       []string{"jira", "create-with-thread"},
	RequireMention: true,
	Callback: func(ctx context.Context, client util.SlackClientInterface, evt *slackevents.MessageEvent, args []string) ([]slack.MsgOption, error) {
		url := util.GetThreadUrl(evt)
		description := ""
		if len(url) > 0 {
			description = fmt.Sprintf("%s\n\ncreated from thread: %s", description, url)
		}
		issue, err := issue.CreateIssue(args[2], "follow up on slack thread", description, args[3])
		if err != nil {
			return util.WrapErrorToBlock(err, "error creating issue"), nil
		}
		issueKey := issue.Key
		issueURL := fmt.Sprintf("%s/browse/%s", JIRA_BASE_URL, issueKey)
		return util.StringToBlock(fmt.Sprintf("issue <%s|%s> created", issueURL, issueKey), false), nil
	},
	RequiredArgs: 4,
	HelpMarkdown: "create a Jira issue with a summary of the thread: `jira create-with-thread [project] [type]`",
	ShouldMatch: []string{
		"jira create-with-thread PROJECT bug",
		"jira create-with-thread PROJECT Todo",
	},
	ShouldntMatch: []string{
		"jira create PROJECT bug",
		"jira create PROJECT Todo",
	},
}
