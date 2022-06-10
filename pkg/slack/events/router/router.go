package router

import (
	"github.com/slack-go/slack"

	"github.com/openshift-eng/splat-sandbox/pkg/slack/events"
	forumvmware "github.com/openshift-eng/splat-sandbox/pkg/slack/events/channels/forum_vmware"
	"github.com/openshift-eng/splat-sandbox/pkg/slack/events/mention"
)

// ForEvents returns a Handler that appropriately routes
// event callbacks for the handlers we know about
func ForEvents(client *slack.Client) events.Handler {
	return events.MultiHandler(
		mention.Handler(client),
		forumvmware.Handler(client),
	)
}
