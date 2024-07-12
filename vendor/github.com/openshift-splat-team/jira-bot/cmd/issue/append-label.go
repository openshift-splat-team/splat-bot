package issue

import (
	"fmt"
	"log"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/openshift-splat-team/jira-bot/pkg/util"
	"github.com/spf13/cobra"
)

func init() {
	cmdAppendLabel.Flags().BoolVarP(&options.dryRunFlag, "dry-run", "d", true, "only apply changes with --dry-run=false")
	cmdAppendLabel.Flags().StringVarP(&options.label, "label", "l", "", "label to append to issues")
	cmdAppendLabel.Flags().StringVarP(&options.skipLabel, "skip-label", "s", "", "if label is present, skip issue")
	err := cmdAppendLabel.MarkFlagRequired("label")
	if err != nil {
		log.Fatalf("unable to mark flag as required: %v", err)
	}
	cmdIssue.AddCommand(cmdAppendLabel)
}

func appendLabel(filter string, options *issueCommandOptions) error {
	client, err := util.GetJiraClient()
	if err != nil {
		return fmt.Errorf("unable to get Jira client: %v", err)
	}

	issues, _, err := client.Issue.Search(filter, &jira.SearchOptions{})
	if err != nil {
		return fmt.Errorf("unable to search for issues: %v", err)
	}

	for _, issue := range issues {
		labels := issue.Fields.Labels
		hasLabel := false
		for _, label := range labels {
			if label == options.label {
				log.Printf("issue: %s already has label: %s", issue.Key, options.label)
				hasLabel = true
				break
			}
			if len(options.skipLabel) > 0 && label == options.skipLabel {
				log.Printf("issue: %s has label: %s, skipping", issue.Key, options.skipLabel)
				hasLabel = true
				break
			}
		}
		if !hasLabel {
			if !options.dryRunFlag {
				_, err := client.Issue.UpdateIssue(issue.Key, map[string]interface{}{
					"fields": map[string]interface{}{
						"labels": append(issue.Fields.Labels, options.label),
					},
				})
				if err != nil {
					return fmt.Errorf("unable to update issue %s: %v", issue.Key, err)
				}
				log.Printf("added label: %s to issue: %s", options.label, issue.Key)
				time.Sleep(5 * time.Second)
			} else {
				log.Printf("issue: %s would have had label: %s added. run again and provide --dry-run=false to apply.", issue.Key, options.label)
			}
		}
	}

	return nil
}

var cmdAppendLabel = &cobra.Command{
	Use:   "append-label [filter]",
	Short: "Appends a label to issues returned from a filter",
	Long:  `Appends a label to issues returned from a filter`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := appendLabel(args[0], &options)
		if err != nil {
			util.RuntimeError(fmt.Errorf("unable to update issue: %v", err))
		}
	},
}
