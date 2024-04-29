package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/openshift-splat-team/splat-bot/data"
	"github.com/openshift-splat-team/splat-bot/pkg/util"
	sbutils "github.com/openshift-splat-team/splat-bot/pkg/util"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

var (
	PROMPT_ISSUE_TITLE   util.Prompt
	PROMPT_ISSUE_SUMMARY util.Prompt
)

func init() {
	val := os.Getenv("PROMPT_ISSUE_TITLE")
	if val != "" {
		PROMPT_ISSUE_TITLE = util.Prompt(val)
	} else {
		PROMPT_ISSUE_TITLE = util.Prompt("can you summarize this thread to a single line? The line length should be less than 100 characters. ")
	}
	val = os.Getenv("PROMPT_ISSUE_SUMMARY")
	if val != "" {
		PROMPT_ISSUE_SUMMARY = util.Prompt(val)
	} else {
		PROMPT_ISSUE_SUMMARY = util.Prompt("can you summarize this thread?")
	}
}

var SummarizeAttributes = data.Attributes{
	Commands:            []string{"summary"},
	RequireMention:      true,
	ResponseIsEphemeral: true,
	Callback: func(ctx context.Context, client util.SlackClientInterface, evt *slackevents.MessageEvent, args []string) ([]slack.MsgOption, error) {
		response, err := util.HandlePrompt(ctx, PROMPT_ISSUE_SUMMARY, client, evt)
		if err != nil {
			return nil, fmt.Errorf("unable to get summary: %v", err)
		}
		return sbutils.StringToBlock(fmt.Sprintf("Sure! Here is a summary of this thread.\n\n*Note: I am a bot and I try my best to provide a reasonable summary. Be sure to check the summary for accuracy.*\n\n%s\n", response), false), nil
	},
	RequiredArgs: 1,
	HelpMarkdown: "summarize this thread: `summary`",
	ShouldMatch: []string{
		"summary",
	},
	ShouldntMatch: []string{
		"jira create-with-summary PROJECT bug",
		"jira create-with-summary PROJECT Todo",
	},
}
