package issue

import (
	"fmt"
	"log"

	"github.com/andygrunwald/go-jira"
	"github.com/openshift-splat-team/jira-bot/pkg/util"
	"github.com/spf13/cobra"
)

// CreateIssue creates an issue in a given project.  the creator of the
// issue will match the user creating the issue.
func CreateIssue(project, summary, description, issueType string) (*jira.Issue, error) {
	return createIssue(&issueCommandOptions{
		summary:     summary,
		description: description,
		issueType:   issueType,
		project:     project,
	})
}

func createIssue(options *issueCommandOptions)  (*jira.Issue, error) {
	client, err := util.GetJiraClient()
	if err != nil {
		return nil, fmt.Errorf("unable to get Jira client: %v", err)
	}

	project, err := util.GetProject(client, options.project)
	if err != nil {
		return nil, fmt.Errorf("unable to get Jira project: %v", err)
	}

	issueType, err := util.GetIssueType(project, options.issueType)
	if err != nil {
		return nil, fmt.Errorf("unable to get Jira issue type: %v", err)
	}

	newIssue := &jira.Issue{
		Fields: &jira.IssueFields{
			Summary:     options.summary,
			Description: options.description,
			Project:     *project,
			Type:        *issueType,
		},
	}

	log.Printf("creating new issue: %+v", newIssue)
	issue, resp, err := client.Issue.Create(newIssue)
	if err != nil {
		responseBody, _ := util.GetResponseBody(resp)
		return nil, fmt.Errorf("unable to create issue: %v. response body: %s", err, responseBody)
	}

	log.Printf("created issue: %s", issue.ID)
	return issue, nil
}

var cmdCreateIssue = &cobra.Command{
	Use:   "create [project] [type] [description] [summary]",
	Short: "creates an issue",
	Long:  `creates an issue`,
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		options.project = args[0]
		options.issueType = args[1]
		options.description = args[2]
		options.summary = args[3]
		_, err := createIssue(&options)
		if err != nil {
			util.RuntimeError(fmt.Errorf("unable to create issue: %v", err))
		}
	},
}
