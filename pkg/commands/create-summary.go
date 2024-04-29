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

var CreateSummaryAttributes = data.Attributes{
	Commands:            []string{"jira", "create-with-summary"},
	RequireMention:      true,
	ResponseIsEphemeral: true,
	Callback: func(ctx context.Context, client util.SlackClientInterface, evt *slackevents.MessageEvent, args []string) ([]slack.MsgOption, error) {
		url := util.GetThreadUrl(evt)
		description := ""
		if len(url) > 0 {
			description = fmt.Sprintf("%s\n\ncreated from thread: %s", description, url)
		}

		response, err := util.HandlePrompt(ctx, PROMPT_ISSUE_SUMMARY, client, evt)
		if err != nil {
			return nil, fmt.Errorf("unable to get summary: %v", err)
		}
		description = fmt.Sprintf("thread summary: %s\n%s\nissue created by splat-bot\n", response, description)

		issue, err := issue.CreateIssue(args[2], "follow up on slack thread", description, args[3])
		if err != nil {
			return util.WrapErrorToBlock(err, "error creating issue"), nil
		}
		issueKey := issue.Key
		issueURL := fmt.Sprintf("%s/browse/%s", JIRA_BASE_URL, issueKey)
		return util.StringToBlock(fmt.Sprintf("issue <%s|%s> created", issueURL, issueKey), false), nil
	},
	RequiredArgs: 4,
	HelpMarkdown: "create a Jira issue with a summary of the thread: `jira create-with-summary [project] [type]`",
	ShouldMatch: []string{
		"jira create-with-summary PROJECT bug",
		"jira create-with-summary PROJECT Todo",
	},
	ShouldntMatch: []string{
		"jira create PROJECT bug",
		"jira create PROJECT Todo",
	},
}
