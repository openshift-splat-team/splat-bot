package commands

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/openshift-splat-team/splat-bot/data"
	"github.com/openshift-splat-team/splat-bot/pkg/util"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

var AskDocsAttributes = data.Attributes{
	Commands:            []string{"ask-docs"},
	RequireMention:      true,
	ResponseIsEphemeral: false,
	AllowNonSplatUsers:  true,
	Callback: func(ctx context.Context, client util.SlackClientInterface, evt *slackevents.MessageEvent, args []string) ([]slack.MsgOption, error) {
		url := os.Getenv("DOC_QUERY_URL")
		if url == "" {
			url = "http://localhost:8000/"
		}

		log.Debugf("question: %v\n", args)
		question := args[1]

		response := "sorry! I was unable to find an answer."

		resp, err := http.Post(url, "text/plain", bytes.NewBufferString(question))
		if err != nil {
			response = fmt.Sprintf("%s. error: %v", response, err)
		}

		if resp.StatusCode == 200 {
			// Read the response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				if err != nil {
					response = fmt.Sprintf("%s. error: %v", response, err)
				}
			} else {
				response = fmt.Sprintf("this doc section may be helpful: <%s|%s>", body, body)
			}
		}
		return util.StringToBlock(response, false), nil
	},
	RequiredArgs: 2,
	HelpMarkdown: "ask docs a question: `ask-docs ask a question`",
	ShouldMatch: []string{
		"ask-docs how do I rotate credentials for vSphere?",
		"ask-docs what version of ESXi is required to install OpenShift?",
	},
	ShouldntMatch: []string{
		"jira create PROJECT bug",
		"jira create PROJECT Todo",
	},
}
