package commands

import (
	"context"
	"fmt"
	"github.com/openshift-splat-team/splat-bot/pkg/controllers"

	"github.com/openshift-splat-team/splat-bot/data"
	"github.com/openshift-splat-team/splat-bot/pkg/util"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func init() {
	AddCommand(PoolsAttributes)
}

var PoolsAttributes = data.Attributes{
	Commands:       []string{"ci", "pools"},
	RequireMention: true,
	Callback: func(ctx context.Context, client util.SlackClientInterface, evt *slackevents.MessageEvent, args []string) ([]slack.MsgOption, error) {
		result := ""
		var err error
		if len(args) > 2 {
			switch args[2] {
			case "schedulable":
				if len(args) < 4 {
					return util.StringToBlock(err.Error(), false), fmt.Errorf("pool schedule requires the name of the pool")
				}
				err = controllers.SetPoolSchedulable(ctx, args[3], true)
				if err != nil {
					return util.StringToBlock(err.Error(), false), fmt.Errorf("failed to set pool schedulable: %w", err)
				}
				result = "pool is schedulable"
			case "unschedulable":
				if len(args) < 4 {
					return util.StringToBlock(err.Error(), false), fmt.Errorf("pool unschedule requires the name of the pool")
				}
				err = controllers.SetPoolSchedulable(ctx, args[3], false)
				if err != nil {
					return util.StringToBlock(err.Error(), false), fmt.Errorf("failed to set pool unschedulable: %w", err)
				}
				result = "pool is unschedulable"
			case "list":
				fallthrough
			default:
				result, err = controllers.GetPoolStatus()
				if err != nil {
					return nil, fmt.Errorf("failed to fetch pool status: %w", err)
				}
			}

		}

		return util.StringToBlock(result, false), nil
	},
	RequiredArgs: 0,
	HelpMarkdown: "interact with vSphere CI pools: `ci pools list|schedulable|unschedulable <pool name>`",
	ShouldMatch: []string{
		"ci pools list",
		"ci pools schedulable pool-1",
		"ci pools unschedulable pool-1",
	},
	ShouldntMatch: []string{
		"jira create-with-summary PROJECT bug",
		"jira create-with-summary PROJECT Todo",
	},
}