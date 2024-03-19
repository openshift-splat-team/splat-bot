package commands

import (
	"context"
	"fmt"

	"github.com/openshift-splat-team/jira-bot/cmd/issue"
	"github.com/openshift-splat-team/splat-bot/data"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

var CreateAttributes = data.Attributes{
	Commands:       []string{"jira", "create"},
	RequireMention: true,
	Callback: func(ctx context.Context, client *socketmode.Client, evt *slackevents.MessageEvent, args []string) ([]slack.MsgOption, error) {
		url := GetThreadUrl(evt)
		description := args[4]
		if len(url) > 0 {
			description = fmt.Sprintf("%s\n\ncreated from thread: %s", description, url)
		}
		description = fmt.Sprintf("%s\nissue created by splat-bot\n", description)

		issue, err := issue.CreateIssue(args[2], args[3], description, args[5])
		if err != nil {
			return WrapErrorToBlock(err, "error creating issue"), nil
		}
		issueKey := issue.Key
		issueURL := fmt.Sprintf("%s/browse/%s", JIRA_BASE_URL, issueKey)
		return StringToBlock(fmt.Sprintf("issue <%s|%s> created", issueURL, issueKey), false), nil
	},
	RequiredArgs: 6,
	HelpMarkdown: "create a Jira issue: `jira create [project] [summary] [description] [type]`",
}
