package mention

import (
	"fmt"
	"strings"

	"github.com/openshift-eng/splat-sandbox/pkg/slack/events"

	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

const (
	channelName = "forum-vmware"
)

var (
	responses = []events.AutoResponseStruct{
		{
			Keywords: []string{"problem", "bug", "issue"},
			Response: "It looks like you are reporting a problem.  To enable us to best help, please gather necessary logs: " +
				"<https://docs.openshift.com/container-platform/4.10/support/gathering-cluster-data.html|must-gather(post install)>, " +
				"<https://docs.openshift.com/container-platform/4.10/support/troubleshooting/troubleshooting-installations.html#installation-bootstrap-gather_troubleshooting-installations|installer-gather(during install)>. " +
				"We'll likely need logs to provide meaningful analysis.  If you already provided them, thanks!!",
			RequiresChannelTag: true,
		},
		{
			Keywords: []string{"migration", "migrate", "moving"},
			Operator: events.RESPONSE_OPERATOR_OR,
			Response: "It looks like you are starting a thread about migration.  Here are some resources which might be helpful: " +
				"<https://access.redhat.com/articles/6718991|Migrating Virtual Machines with vMotion>, " +
				"<https://docs.openshift.com/container-platform/4.10/migration_toolkit_for_containers/about-mtc.html#migration-direct-volume-migration-and-direct-image-migration_about-mtc|Migrating Persistent Volumes with Direct Volume Migration>. " +
				"Someone will follow up on this thread when able.",
			RequiresChannelTag: true,
		},
		{
			Keywords: []string{"performance"},
			Operator: events.RESPONSE_OPERATOR_OR,
			Response: "It looks like you are starting a thread about performance.  Here are some resources which might be helpful: " +
				"<https://access.redhat.com/articles/5822821|Triaging OpenShift Performance on VMware>, " +
				"<https://communities.vmware.com/t5/Storage-Performance/Interpreting-esxtop-Statistics/ta-p/2776936|Interpreting esxtop statistics>. " +
				"Someone will follow up on this thread when able.",
			RequiresChannelTag: true,
		},
		{
			Keywords: []string{"SRM", "supported"},
			Operator: events.RESPONSE_OPERATOR_AND,
			Response: "It looks like you are starting a thread about SRM.  SRM is not supported with OpenShift at this time.  This topic has been previously raised in this channel if you'd like to checkout " +
				"previous discussions.",
			RequiresChannelTag: true,
		},
		{
			Keywords: []string{"open-vm-tools", "version"},
			Operator: events.RESPONSE_OPERATOR_AND,
			Response: "It looks like you are starting a thread about the version of open-vm-tools.  The version of open-vm-tools is not upgradable in RHCOS.  " +
				"This is a topic that has been been previously discussed in this channel if you'd like to peruse prior discussions. #forum-coreos may be able to " +
				"provide additional context.",
			RequiresChannelTag: true,
		},
		{
			Keywords: []string{"splat"},
			Operator: events.RESPONSE_OPERATOR_OR,
			Response: "Hey! You mentioned SPLAT.  If it's urgent you can message @splat-team and we'll respond if able. " +
				"If there is something you'd like us to research or follow up on, feel free to create a card on our <https://issues.redhat.com/secure/RapidBoard.jspa?projectKey=SPLAT&rapidView=5962|board>.",
			RequiresChannelTag: false,
		},
	}
)

type messagePoster interface {
	PostMessage(channelID string, options ...slack.MsgOption) (string, string, error)
}

// Handler returns a handler that knows how to respond to
// new messages that mention the robot by showing users
// which interactive workflows they might be interested in,
// based on the phrasing that they used to mention the bot.
func Handler(client messagePoster) events.PartialHandler {
	return events.PartialHandlerFunc("message", func(callback *slackevents.EventsAPIEvent, logger *logrus.Entry) (handled bool, err error) {
		if callback.Type != slackevents.CallbackEvent {
			return false, nil
		}
		event, ok := callback.InnerEvent.Data.(*slackevents.MessageEvent)

		if !ok {
			return false, nil
		}
		if event.BotID != "" {
			logger.Debug("event came from a bot, ignoring")
			return false, nil
		}
		logger.Debug("Handling #forum-vmware message...")
		timestamp := event.TimeStamp

		if event.ThreadTimeStamp != "" {
			timestamp = event.ThreadTimeStamp
		}

		channelId := fmt.Sprintf("<#%s|%s>", event.Channel, channelName)
		channelMatch := strings.Contains(event.Text, channelId)
		blocks := events.ResponseFor(event.Text, channelMatch, responses)
		if len(blocks) == 0 {
			return false, nil
		}
		responseChannel, responseTimestamp, err := client.PostMessage(event.Channel, slack.MsgOptionBlocks(blocks...), slack.MsgOptionTS(timestamp))
		if err != nil {
			logger.WithError(err).Warn("Failed to post response to app mention")
		} else {
			logger.Infof("Posted response to app mention in channel %s at %s", responseChannel, responseTimestamp)
		}
		return true, err
	})
}
