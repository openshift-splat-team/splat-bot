package llm

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"

	"github.com/openshift-splat-team/jira-bot/pkg/util"
)

const (
	PROMPT_RESPONSE_TIMEOUT = time.Second * 120
)

func GetJiraIssueSummary(ctx context.Context, query string) (string, error) {
	client, err := util.GetJiraClient()
	if err != nil {
		return "", fmt.Errorf("unable to get Jira client: %v", err)
	}

	issues, resp, err := client.Issue.Search("filter = \"SPLAT - updates in last week\"", nil) //util.GetIssuesInQuery(client, query)
	if err != nil {
		responseBody, _ := util.GetResponseBody(resp)
		return "", fmt.Errorf("unable to get issues: %v\n\n%s", err, responseBody)
	}

	builder := strings.Builder{}
	for _, issue := range issues {
		log.Printf("issue: %s", issue.Fields.Summary)
		builder.WriteString(fmt.Sprintf("%s: %s\n", issue.Fields.Summary, issue.Fields.Status.ID))
	}
	return builder.String(), nil
}

// GenerateResponse generates a response from an ollama API endpoint
func GenerateResponse(ctx context.Context, prompt string) (string, error) {
	/*endpoint := os.Getenv("OLLAMA_ENDPOINT")
	if len(endpoint) == 0 {
		return "", errors.New("OLLAMA_ENDPOINT must be exported")
	}

	model := os.Getenv("OLLAMA_MODEL")
	if len(model) == 0 {
		model = "tinyllama"
	}

	llm, err := ollama.New(ollama.WithModel(model), ollama.WithServerURL(endpoint))
	if err != nil {
		log.Fatal(err)
	}

	timedCtx, cancel := context.WithTimeout(ctx, PROMPT_RESPONSE_TIMEOUT)
	defer cancel()
	completion, err := llms.GenerateFromSinglePrompt(timedCtx, llm, prompt)
	if err != nil {
		log.Fatal(err)
	}
	return completion, nil*/

	//summay   := GetJiraIssueSummary(ctx, "SPLAT - updates in last wee
	out, err := Completion(prompt)

	if err != nil {
		return "", fmt.Errorf("unable to get completion: %v", err)
	}
	if len(out.Choices) == 0 {
		return "", fmt.Errorf("no completion returned")
	}

	return out.Choices[0].Message.Content, nil
}
