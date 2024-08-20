package issue

import (
	"fmt"
	"log"

	"github.com/andygrunwald/go-jira"
	"github.com/openshift-splat-team/jira-bot/pkg/util"
	"github.com/spf13/cobra"
)

func checkSetSpikePoints(client *jira.Client, issue jira.Issue, options *issueCommandOptions) error {
	if issue.Fields.Type.Name == "Spike" {
		if util.GetStoryPoints(issue.Fields.Unknowns) > 0 && !options.overrideFlag {
			log.Fatalf("issue: %s already has assigned story points.  run again and provide --override=true to apply", issue.Key)
			return nil
		}

		if util.GetStoryPoints(issue.Fields.Unknowns) == 0 || options.overrideFlag {
			propertyMap := map[string]interface{}{
				"fields": map[string]interface{}{
					util.FieldStoryPoints: options.defaultSpikeStoryPoints,
				},
			}
			if options.dryRunFlag {
				log.Printf("issue: %s would have default spike points assigned. run again and provide --dry-run=false to apply.", issue.Key)
				return nil
			} else {
				log.Printf("setting default story points for spike: %s", issue.Key)
				_, err := client.Issue.UpdateIssue(issue.Key, propertyMap)
				if err != nil {
					return fmt.Errorf("unable to update issue %s: %v", issue.Key, err)
				}
			}
		}
	}
	return nil
}

// autoUpdateIssuesInQuery according to rules set forth by the team
func autoUpdateIssuesInQuery(jql string, options *issueCommandOptions) error {
	log.Printf("preparing to auto-update issues found in query: %s", jql)
	jiraClient, err := util.GetJiraClient()
	if err != nil {
		return fmt.Errorf("unable to get Jira client: %v", err)
	}

	issues, _, err := jiraClient.Issue.Search(jql, nil)
	if err != nil {
		return fmt.Errorf("unable to get issues: %v", err)
	}

	log.Printf("%d issues found in query", len(issues))

	for _, issue := range issues {
		if options.defaultSpikeStoryPoints > 0 {
			err = checkSetSpikePoints(jiraClient, issue, options)
			if err != nil {
				return fmt.Errorf("unable to set default spike story points: %v", err)
			}
		}
	}
	return nil
}

var cmdAutoUpdateIssuesStatus = &cobra.Command{
	Use:   "auto-update-issues [jql]",
	Short: "Updates issues according to rules provided as options.",
	Long:  `Updates issues matching the JQL provided according to rules provided as options`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := autoUpdateIssuesInQuery(args[0], &options)
		if err != nil {
			util.RuntimeError(fmt.Errorf("unable to update issues: %v", err))
		}
	},
}
