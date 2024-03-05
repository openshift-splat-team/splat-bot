package util

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/spf13/viper"
)

const (
	FieldStoryPoints   = "customfield_12310243"
	FieldStatusSummary = "customfield_12320841"
)

func GetJiraClient() (*jira.Client, error) {
	token := viper.GetString("personal_access_token")

	tp := jira.BearerAuthTransport{
		Token: token,
	}

	return jira.NewClient(tp.Client(), "https://issues.redhat.com/")
}

func GetIssuesInQuery(client *jira.Client, query string) ([]jira.Issue, []string, error) {
	log.Printf("invoking query: %s\n", query)
	issues, _, err := client.Issue.Search(query, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to execute query: %v", err)
	}
	issueIds := []string{}
	for _, issue := range issues {
		issueIds = append(issueIds, issue.ID)
	}
	log.Printf("found %d issues\n", len(issues))
	return issues, issueIds, nil
}

func GetStoryPoints(totalMap map[string]interface{}) float64 {
	if points, exists := totalMap[FieldStoryPoints]; exists {
		if points != nil {
			return points.(float64)
		}
	}
	return 0
}

// GetIssueType retrieves the identified issue type from Jira
func GetIssueType(project *jira.Project, typeID string) (*jira.IssueType, error) {
	log.Printf("getting issue type: %s", typeID)

	for _, issueType := range project.IssueTypes {
		if strings.EqualFold(issueType.Name, typeID) {
			return &issueType, nil
		}
	}

	return nil, fmt.Errorf("unable to find issue type: %s", typeID)
}

// GetUser retrieves the identified user from Jira
func GetUser(client *jira.Client, userID string) (*jira.User, error) {
	log.Printf("getting user: %s", userID)

	user, _, err := client.User.GetSelf()
	if err != nil {
		return nil, fmt.Errorf("unable to get user: %v", err)
	}
	return user, nil
}

// GetProject retrieves the identified project from Jira
func GetProject(client *jira.Client, projectID string) (*jira.Project, error) {
	log.Printf("getting project: %s", projectID)
	project, _, err := client.Project.Get(projectID)
	if err != nil {
		return nil, fmt.Errorf("unable to get project: %v", err)
	}
	return project, nil
}

func GetResponseBody(resp *jira.Response) (string, error) {
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read response body: %v", err)

	}
	return string(body), nil
}
