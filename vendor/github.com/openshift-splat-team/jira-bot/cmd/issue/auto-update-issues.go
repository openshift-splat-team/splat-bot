package issue

import (
	"fmt"
	"log"
	"strings"
	"time"

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

func getFixVersionsFromParent(client *jira.Client, parentIssueKey string, issue jira.Issue) ([]*jira.FixVersion, bool, error) {
	noFixVersion := false
	log.Printf("getting fix versions from parent: %s", parentIssueKey)
	parentIssue, _, err := client.Issue.Get(parentIssueKey, nil)
	if err != nil {
		return nil, noFixVersion, fmt.Errorf("unable to get parent issue: %v", err)
	}
	for _, fixVersion := range parentIssue.Fields.Labels {
		if strings.Contains(fixVersion, "splat-nofixversion") {
			noFixVersion = true
			break
		}
	}
	return parentIssue.Fields.FixVersions, noFixVersion, err
}

func checkSetFixVersionFromParent(client *jira.Client, issue jira.Issue) error {
	epic, feature := util.GetParentLinks(issue.Fields.Unknowns)
	var fixVersions []*jira.FixVersion
	var err error
	var noFixVersion bool

	if len(epic) > 0 {
		fixVersions, noFixVersion, err = getFixVersionsFromParent(client, epic, issue)
		if err != nil {
			log.Printf("unable to get fix version from epic: %v", err)
		}
	}
	if len(fixVersions) == 0 && len(feature) > 0 {
		fixVersions, noFixVersion, err = getFixVersionsFromParent(client, feature, issue)
		if err != nil {
			log.Printf("unable to get fix version from feature: %v", err)
		}
	}

	labels := issue.Fields.Labels

	if len(fixVersions) == 0 && !noFixVersion {
		log.Printf("no fix version found in parent")
		return nil
	} else {
		if noFixVersion {
			needToAdd := true
			for _, label := range labels {
				if strings.Contains(label, "splat-nofixversion") {
					needToAdd = false
					break
				}
			}
			if needToAdd {
				labels = append(labels, "splat-nofixversion")
				log.Printf("splat-nofixversion label found in parent, adding splat-nofixversion label to issue")
			}
		}
	}

	propertyMap := map[string]interface{}{
		"fields": map[string]interface{}{
			"fixVersions": fixVersions,
			"labels":      labels,
		},
	}
	versions := []string{}
	for _, fixVersion := range fixVersions {
		versions = append(versions, fixVersion.Name)
	}

	log.Printf("setting fix version for issue: %s to %s", issue.ID, strings.Join(versions, ","))
	if options.dryRunFlag {
		log.Printf("issue: %s would have had its fixVersions set.", issue.Key)
		return nil
	} else {
		_, err = client.Issue.UpdateIssue(issue.Key, propertyMap)
		if err != nil {
			return fmt.Errorf("unable to update issue %s: %v", issue.Key, err)
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
		if len(issue.Fields.FixVersions) == 0 {
			time.Sleep(5 * time.Second)
			err = checkSetFixVersionFromParent(jiraClient, issue)
			if err != nil {
				return fmt.Errorf("unable to get fix version from parent: %v", err)
			}
		}

		if options.defaultSpikeStoryPoints > 0 {
			time.Sleep(5 * time.Second)
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
