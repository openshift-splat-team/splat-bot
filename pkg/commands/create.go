package commands

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"text/template"

	log "github.com/sirupsen/logrus"

	"github.com/openshift-splat-team/jira-bot/cmd/issue"
	"github.com/openshift-splat-team/splat-bot/data"
	"github.com/openshift-splat-team/splat-bot/pkg/util"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

const issueTemplateSource = `
*User Story:*
As an {{.Principal}} I want {{.Goal}} so {{.Outcome}}.

*Description:*
< Record any background information >

*Acceptance Criteria:*
< Record how we'll know we're done >

*Other Information:*
< Record anything else that may be helpful to someone else picking up the card >

issue created by splat-bot
`

var issueTemplate *template.Template

func init() {
	var err error

	issueTemplate, err = template.New("issue").Parse(issueTemplateSource)
	if err != nil {
		panic(err)
	}
}

type assistantContext struct {
	Principal string
	Goal      string
	Outcome   string
}

func invokeTemplate(t *template.Template, data any) (string, error) {
	// Create a bytes.Buffer to hold the output
	var buf bytes.Buffer

	// Execute the template and write the result into the buffer
	err := t.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

var CreateAttributes = data.Attributes{
	Commands:       []string{"jira", "create"},
	RequireMention: true,
	Callback: func(ctx context.Context, client util.SlackClientInterface, evt *slackevents.MessageEvent, args []string) ([]slack.MsgOption, error) {
		var description string
		var err error

		fmt.Fprintln(os.Stderr, "processing Jira Create")
		fmt.Fprintf(os.Stderr, "evt: %+v", evt)
		assistantCtx := assistantContext{
			Principal: "OpenShift Engineer",
			Goal:      "___",
			Outcome:   "___",
		}
		url := util.GetThreadUrl(evt)
		log.Debugf("%v", args)
		summary := args[2]

		if len(args) >= 4 {
			assistantCtx.Goal = args[2]
			assistantCtx.Outcome = args[3]
		}

		// Execute the template and write the result into the buffer
		description, err = invokeTemplate(issueTemplate, assistantCtx)
		if err != nil {
			return util.StringToBlock(fmt.Sprintf("unable to to process template. error: %v", err), false), nil
		}

		if len(url) > 0 {
			description = fmt.Sprintf("%s\n\ncreated from thread: %s", description, url)
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
	MaxArgs:      4,
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
