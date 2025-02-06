package commands

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"

	"github.com/openshift-splat-team/splat-bot/data"
	"github.com/openshift-splat-team/splat-bot/pkg/controllers"
	"github.com/openshift-splat-team/splat-bot/pkg/util"
)

func init() {
	AddCommand(LeasesAttributes)
}

type leaseOptions struct {
	cpus     int
	memory   int
	networks int
	pool     string
}

func getLeaseOptions(args []string) leaseOptions {
	cpus := 24
	memory := 96
	networks := 1
	pools := ""
	log.Printf("lease args: %v", args)
	if len(args) >= 4 {
		log.Printf("applying options to lease")
		for _, arg := range args[3:] {
			parts := strings.Split(arg, "=")
			if len(parts) != 2 {
				continue
			}
			switch parts[0] {
			case "cpus":
				cpus, _ = strconv.Atoi(parts[1])
			case "memory":
				memory, _ = strconv.Atoi(parts[1])
			case "networks":
				networks, _ = strconv.Atoi(parts[1])
			case "pools":
				// Need to remove the double quotes if added for multiple pools
				pools = strings.Replace(parts[1], "\"", "", -1)
			}
		}
	}
	return leaseOptions{
		cpus:     cpus,
		memory:   memory,
		networks: networks,
		pool:     pools,
	}
}

func validateLeaseOptions(ctx context.Context, options leaseOptions) error {
	// Validate the pool names.  An incorrect pool name will lead to a bad time
	if len(options.pool) > 0 {
		pools, err := controllers.GetPoolNames(ctx)
		if err != nil {
			return err
		}

		found := false
		for _, curPool := range pools {
			if curPool == options.pool {
				found = true
				break
			}
		}

		if !found {
			return errors.New("pool not found")
		}
	}
	return nil
}

var LeasesAttributes = data.Attributes{
	Commands:       []string{"ci", "lease"},
	RequireMention: true,
	Callback: func(ctx context.Context, client util.SlackClientInterface, evt *slackevents.MessageEvent, args []string) ([]slack.MsgOption, error) {
		result := ""
		var err error
		if len(args) > 2 {
			switch args[2] {
			case "acquire":
				options := getLeaseOptions(args)

				if err = validateLeaseOptions(ctx, options); err != nil {
					return util.StringToBlock(err.Error(), false), fmt.Errorf("failed to acquire lease: %w", err)
				}

				_, err := controllers.AcquireLease(ctx, evt.User, options.cpus, options.memory, options.pool, options.networks)
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
		"ci lease list",
		"ci lease acquire (optional args) cpus=24 memory=96 networks=1 pools=\"space-separated-pool-names\"",
		"ci lease release",
	},
	ShouldntMatch: []string{
		"jira create-with-summary PROJECT bug",
		"jira create-with-summary PROJECT Todo",
	},
}
