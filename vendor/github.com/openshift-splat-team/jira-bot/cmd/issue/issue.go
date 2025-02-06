package issue

import (
	"github.com/spf13/cobra"
)

type issueCommandOptions struct {
	defaultSpikeStoryPoints int64
	points                  int64
	dryRunFlag              bool
	overrideFlag            bool
	state                   string
	priority                string
	summary                 string
	label                   string
	skipLabel               string
	description             string
	issueType               string
	project                 string
}

var options = issueCommandOptions{
	defaultSpikeStoryPoints: -1,
	dryRunFlag:              true,
	overrideFlag:            false,
	priority:                "",
}

var cmdIssue = &cobra.Command{
	Use:   "issue",
	Short: "Manage issues",
	Long:  `This command allows you to manage issues in your project management tool.`,
}

func Initialize(rootCmd *cobra.Command) {
	cmdUpdateSizeAndPriority.Flags().BoolVarP(&options.dryRunFlag, "dry-run", "d", true, "only apply changes with --dry-run=false")
	cmdUpdateSizeAndPriority.Flags().BoolVarP(&options.overrideFlag, "override", "o", false, "overrides a warning when --override=true")
	cmdUpdateSizeAndPriority.Flags().Int64VarP(&options.points, "points", "p", -1, "points to apply to issue")
	cmdUpdateSizeAndPriority.Flags().StringVarP(&options.priority, "priority", "r", "", "priority to set")
	cmdUpdateSizeAndPriority.Flags().StringVarP(&options.state, "state", "s", "", "sets the issue state")

	cmdTriageIssues.Flags().BoolVarP(&options.dryRunFlag, "dry-run", "d", true, "only apply changes with --dry-run=false")

	cmdAutoUpdateIssuesStatus.Flags().BoolVarP(&options.overrideFlag, "override", "o", false, "overrides a warning when --override=true")
	cmdAutoUpdateIssuesStatus.Flags().Int64VarP(&options.defaultSpikeStoryPoints, "default-spike-points", "s", -1, "points to apply to spikes which have no points")
	cmdAutoUpdateIssuesStatus.Flags().BoolVarP(&options.dryRunFlag, "dry-run", "d", true, "only apply changes with --dry-run=false")

	cmdIssue.AddCommand(cmdCreateIssue)
	cmdIssue.AddCommand(cmdAutoUpdateIssuesStatus)
	cmdIssue.AddCommand(cmdUpdateSizeAndPriority)
	rootCmd.AddCommand(cmdIssue)
}
