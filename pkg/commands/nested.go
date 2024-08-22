package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"

	"github.com/openshift-splat-team/splat-bot/data"
	"github.com/openshift-splat-team/splat-bot/pkg/controllers"
	"github.com/openshift-splat-team/splat-bot/pkg/util"
)

const (
	repoUrl  = "https://github.com/"
	repoPath = "/tmp/repo"
)

var Nested = data.Attributes{
	Commands:       []string{"ci", "nested", "lease"},
	RequireMention: true,
	Callback: func(ctx context.Context, client util.SlackClientInterface, evt *slackevents.MessageEvent, args []string) ([]slack.MsgOption, error) {
		result := ""
		var err error
		if len(args) > 2 {
			switch args[2] {
			case "acquire":
				options := getLeaseOptions(args)

				if _, err := os.Stat(repoPath); os.IsNotExist(err) {
					if err = os.MkdirAll(repoPath, 0755); err != nil {
						return nil, err
					}

					err := util.Clone(repoUrl, repoPath)
					if err != nil {
						return util.StringToBlock(err.Error(), false), fmt.Errorf("failed to clone repo: %w", err)
					}
				}

				_, err = controllers.AcquireLease(ctx, evt.User, options.cpus, options.memory, options.pool, options.networks, controllers.SplatBotNestedLease)
				if err != nil {
					return util.StringToBlock(err.Error(), false), fmt.Errorf("failed to acquire lease: %w", err)
				}
				result = "Lease(s) have been created. Once fulfilled by the vSphere capacity manager you will receive a direct message " +
					"with further details. This could take a few minutes."

			case "renew":
				expires, err := controllers.RenewLease(ctx, evt.User)
				if err != nil {
					return util.StringToBlock(err.Error(), false), fmt.Errorf("failed to renew lease: %w", err)
				}
				result = fmt.Sprintf("Your lease has been renewed. It expires at %s", expires)
			case "release":
				err = controllers.RemoveLease(ctx, evt.User)
				if err != nil {
					return util.StringToBlock(err.Error(), false), fmt.Errorf("failed to set pool unschedulable: %w", err)
				}
				result = "Your lease(s) and associated resources are being deleted. You will receive a notification when this is complete."
			case "list":
				fallthrough
			default:
				result, err = controllers.GetLeaseStatus(evt.User)
				if err != nil {
					return util.StringToBlock(err.Error(), false), fmt.Errorf("failed to fetch pool status: %w", err)
				}
			}
		}

		return util.StringToBlock(result, false), nil
	},
	RequiredArgs: 0,
	HelpMarkdown: "interact with your vSphere CI leases: `ci lease list|acquire|release`",
	ShouldMatch: []string{
		"ci nested lease list",
		"ci nested lease acquire (optional args) cpus=24 memory=96 networks=1 pools=\"space-separated-pool-names\"",
		"ci nested lease release",
	},
	ShouldntMatch: []string{
		"jira create-with-summary PROJECT bug",
		"jira create-with-summary PROJECT Todo",
	},
}
