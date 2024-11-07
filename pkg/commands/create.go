package commands

import (
	"bytes"
	"context"
	"fmt"
	"github.com/openshift-splat-team/jira-bot/cmd/issue"
	"log"
	"strings"
	"text/template"

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

const assistantTemplateSource = `You are an engineer who writes user stories. You take a principal, a goal, and a desired outcome to create a user story with a story, a description, and acceptance criteria. The principal is {{.Principal}}, the goal is {{.Goal}}, and the desired outcome is {{.Outcome}}.`

const assistantTemplateResponse = `You provided enough information to attempt to generate a sample story based on the information provided.
This sample story should provide a helpful starting point for your card. However, this is not a perfect system and you'll
need to review and likely refine the sample to make sense.


{{.Sample}}`

var issueTemplate *template.Template
var assistantTemplate *template.Template
var assistantResponse *template.Template

var issueTypeMap = map[string]string{
	"task":  "Task",
	"bug":   "Bug",
	"spike": "Spike",
	"story": "Story",
}

func init() {
	var err error

	issueTemplate, err = template.New("issue").Parse(issueTemplateSource)
	if err != nil {
		panic(err)
	}

	assistantTemplate, err = template.New("assistant").Parse(assistantTemplateSource)
	if err != nil {
		panic(err)
	}

	assistantResponse, err = template.New("assistantResponse").Parse(assistantTemplateResponse)
	if err != nil {
		panic(err)
	}
}

type assistantContext struct {
	Principal string
	Goal      string
	Outcome   string
	Type      string
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

		assistantCtx := assistantContext{
			Principal: "OpenShift Engineer",
			Goal:      "___",
			Outcome:   "___",
			Type:      "Task",
		}
		url := util.GetThreadUrl(evt)
		fmt.Printf("%v", args)
		summary := args[2]
		outcome := ""

		if len(args) >= 4 {
			outcome = args[3]
			assistantCtx.Goal = args[2]
			assistantCtx.Outcome = args[3]
		}

		if len(args) >= 5 {
			issueType := args[4]
			if _it, exists := issueTypeMap[strings.ToLower(issueType)]; exists {
				assistantCtx.Type = _it
			} else {
				return util.StringToBlock("supported issue types are task(default), story, bug, spike", false), nil
			}
		}
		// Execute the template and write the result into the buffer
		description, err = invokeTemplate(issueTemplate, assistantCtx)
		if err != nil {
			return util.StringToBlock(fmt.Sprintf("unable to to process template. error: %v", err), false), nil
		}

		if len(url) > 0 {
			description = fmt.Sprintf("%s\n\ncreated from thread: %s", description, url)
		}

		issue, err := issue.CreateIssue("SPLAT", summary, description, assistantCtx.Type)
		if err != nil {
			return util.WrapErrorToBlock(err, "error creating issue"), nil
		}
		issueKey := issue.Key
		issueURL := fmt.Sprintf("%s/browse/%s", JIRA_BASE_URL, issueKey)

		if len(outcome) > 0 {
			log.Print("requesting sample issue from LLM")
			prompt, err := invokeTemplate(assistantTemplate, assistantCtx)
			if err != nil {
				return util.StringToBlock(fmt.Sprintf("unable to to process request template. error: %v", err), false), nil
			}

			handlePrompt, err := util.GenerateResponse(ctx, prompt)
			if err != nil {
				return util.StringToBlock(fmt.Sprintf("unable to handle assistant template. error: %v", err), false), nil
			}

			params := map[string]string{
				"Sample": handlePrompt,
			}
			response, err := invokeTemplate(assistantResponse, params)
			if err != nil {
				return util.StringToBlock(fmt.Sprintf("unable to to process template. error: %v", err), false), nil
			}

			socketClient, err := util.GetClient()
			if err != nil {
				return nil, fmt.Errorf("unable to to process template. error: %v", err)
			}

			msgOptions := util.StringToBlock(response, false)

			_, err = socketClient.PostEphemeral(evt.Channel, evt.User, msgOptions...)
			if err != nil {
				return nil, fmt.Errorf("unable to to process template. error: %v", err)
			}
		}
		return util.StringToBlock(fmt.Sprintf("issue <%s|%s> created", issueURL, issueKey), false), nil
	},
	RequiredArgs: 3,
	MaxArgs:      5,
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
