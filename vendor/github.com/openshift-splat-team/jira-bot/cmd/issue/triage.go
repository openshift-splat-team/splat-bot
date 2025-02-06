package issue

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/openshift-splat-team/jira-bot/pkg/util"
	"github.com/spf13/cobra"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

func init() {
	cmdIssue.AddCommand(cmdTriageIssues)
}

var re = regexp.MustCompile(`[\r\n,\t]`)

func applyDeleteLabel(client *jira.Client, issue jira.Issue, label string, del bool) error {
	labels := issue.Fields.Labels
	updatedIssue, _, err := client.Issue.Get(issue.Key, nil)
	if err != nil {
		return fmt.Errorf("failed to get issue by key %s: %v", issue.Key, err)
	}
	issue = *updatedIssue
	hasLabel := false
	for _, _label := range labels {
		if _label == label {
			if !del {
				log.Printf("issue: %s already has label: %s", issue.Key, label)
			}
			hasLabel = true
			break
		}
	}
	if !hasLabel {
		action := "added"
		if del {
			action = "removed"
		}
		if !options.dryRunFlag {
			updatedLabels := issue.Fields.Labels
			if del {
				updatedLabels = []string{}
				for _, _label := range issue.Fields.Labels {
					if _label != label {
						updatedLabels = append(updatedLabels, _label)
					}
				}
			} else {
				updatedLabels = append(updatedLabels, label)
			}
			_, err := client.Issue.UpdateIssue(issue.Key, map[string]interface{}{
				"fields": map[string]interface{}{
					"labels": updatedLabels,
				},
			})
			if err != nil {
				return fmt.Errorf("unable to update issue %s: %v", issue.Key, err)
			}
			log.Printf("%s label: %s to issue: %s", action, label, issue.Key)
			time.Sleep(5 * time.Second)
		} else {
			log.Printf("issue: %s would have had label: %s %s. run again and provide --dry-run=false to apply.", issue.Key, label, action)
		}
	}
	return nil
}
func getTriageBuckets(issues []jira.Issue) (map[string][]string, map[string][]string) {
	storyPointTriageMap := map[string][]string{"Description": make([]string, 0), "Summary": make([]string, 0), "Issue Type": make([]string, 0), "Issue Key": make([]string, 0)}
	readyForRefinementTriageMap := map[string][]string{"Description": make([]string, 0), "Summary": make([]string, 0), "Issue Key": make([]string, 0)}

	for _, issue := range issues {
		log.Printf("issue %s", issue.Key)
		if issue.Fields.Type.Name == "Bug" {
			fmt.Printf("Issue %s is a bug, bugs don't get triaged.  skipping...\n", issue.Key)
			continue
		}
		description := issue.Fields.Description
		description = re.ReplaceAllString(description, " ")
		summary := issue.Fields.Summary
		summary = re.ReplaceAllString(summary, " ")
		log.Printf("check issue %s for story points", issue.Key)
		if util.GetStoryPoints(issue.Fields.Unknowns) == 0 {
			storyPointTriageMap["Description"] = append(storyPointTriageMap["Description"], description)
			storyPointTriageMap["Summary"] = append(storyPointTriageMap["Summary"], summary)
			storyPointTriageMap["Issue Type"] = append(storyPointTriageMap["Issue Type"], issue.Fields.Type.Name)
			storyPointTriageMap["Issue Key"] = append(storyPointTriageMap["Issue Key"], issue.Key)
		}

		checkRefinement := true
		log.Printf("labels %v for issue %s", issue.Fields.Labels, issue.Key)
		for _, label := range issue.Fields.Labels {
			switch label {
			case "ready-for-prioritization":
				log.Printf("issue %s already has ready-for-prioritization", issue.Key)
				checkRefinement = false
			}
		}
		if checkRefinement {
			readyForRefinementTriageMap["Description"] = append(readyForRefinementTriageMap["Description"], description)
			readyForRefinementTriageMap["Summary"] = append(readyForRefinementTriageMap["Summary"], summary)
			readyForRefinementTriageMap["Issue Key"] = append(readyForRefinementTriageMap["Issue Key"], issue.Key)
		}
	}

	return storyPointTriageMap, readyForRefinementTriageMap
}

// ;./jira-bot issue generate-sizings "filter = \"OpenShift SPLAT - No story points assigned\"" --dry-run=false
func generateIssueSizings(client *jira.Client, issues []jira.Issue, storyPointTriageMap map[string][]string, options *issueCommandOptions) error {
	if len(storyPointTriageMap["Description"]) == 0 {
		log.Print("no issues found which need story points")
		return nil
	}
	url := os.Getenv("JIRA_NEURAL_SIZING_URL")
	if len(url) == 0 {
		url = "http://127.0.0.1:8001"
	}
	outMap, err := util.PostJSONData(url, storyPointTriageMap)
	if err != nil {
		return fmt.Errorf("error sending request to %s: %v", url, err)
	}

	for _, issue := range issues {
		foundIssue := false
		for k, v := range outMap["Issue Key"] {
			if v != issue.Key {
				continue
			}
			val, err := strconv.Atoi(outMap["sizing"][k])
			if err != nil {
				fmt.Printf("Issue %s is not a valid size.  skipping...\n", issue.Key)
				continue
			}
			foundIssue = true
			options.points = int64(val)
			break
		}
		if !foundIssue {
			continue
		}

		updated, err := checkSetPoints(client, issue, options)
		if err != nil {
			return fmt.Errorf("error setting points on issue: %v", err)
		}

		if updated {
			_, _, err = client.Issue.AddComment(issue.Key, &jira.Comment{
				Body: "set sizing based on past history of issues sized by the team. this issue should still be sized and refined.",
				Visibility: jira.CommentVisibility{
					Type:  "role",
					Value: "Administrators",
				}})
			if err != nil {
				return fmt.Errorf("unable to add comment %s: %v", url, err)
			}
		}
	}
	return nil
}

func triageIssueRefinement(client *jira.Client, issues []jira.Issue, refinementTriageMap map[string][]string, options *issueCommandOptions) error {
	if len(refinementTriageMap["Description"]) == 0 {
		log.Print("no issues found which need refinement triage")
		return nil
	}

	url := os.Getenv("JIRA_NEURAL_REFINEMENT_URL")
	if len(url) == 0 {
		url = "http://127.0.0.1:8002"
	}
	outMap, err := util.PostJSONData(url, refinementTriageMap)
	if err != nil {
		return fmt.Errorf("error sending request to %s: %v", url, err)
	}

	for _, issue := range issues {
		responseIdx := ""
		for k, v := range outMap["Issue Key"] {
			if v != issue.Key {
				continue
			}
			responseIdx = k
			break
		}
		if len(responseIdx) == 0 {
			continue
		}
		if outMap["REFINED"][responseIdx] == "true" {
			err = applyDeleteLabel(client, issue, "ready-for-prioritization", true)
			if err != nil {
				return fmt.Errorf("error applying label: %v", err)
			}
			err = applyDeleteLabel(client, issue, "needs-refinement", false)
			if err != nil {
				return fmt.Errorf("error applying label: %v", err)
			}
		} else {
			err = applyDeleteLabel(client, issue, "ready-for-prioritization", false)
			if err != nil {
				return fmt.Errorf("error applying label: %v", err)
			}
			err = applyDeleteLabel(client, issue, "needs-refinement", true)
			if err != nil {
				return fmt.Errorf("error applying label: %v", err)
			}

		}
	}
	return nil
}

func triageIssues(filter string, options *issueCommandOptions) error {
	client, err := util.GetJiraClient()
	if err != nil {
		return fmt.Errorf("error getting jira client: %v", err)
	}
	issues, _, err := client.Issue.Search(filter, &jira.SearchOptions{})
	if err != nil {
		return fmt.Errorf("unable to search for issues: %v", err)
	}

	pointsMap, refinedMap := getTriageBuckets(issues)
	err = generateIssueSizings(client, issues, pointsMap, options)
	if err != nil {
		return fmt.Errorf("error generating sizings: %v", err)
	}
	err = triageIssueRefinement(client, issues, refinedMap, options)
	if err != nil {
		return fmt.Errorf("error triaging refinements: %v", err)
	}
	return nil
}

var cmdTriageIssues = &cobra.Command{
	Use:   "triage-issues [filter]",
	Short: "Triages issues",
	Long:  `Triages issues`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := triageIssues(args[0], &options)
		if err != nil {
			util.RuntimeError(fmt.Errorf("unable to triage issues: %v", err))
		}
	},
}
