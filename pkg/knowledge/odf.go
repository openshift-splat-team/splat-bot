package knowledge

import (
	"github.com/openshift-splat-team/splat-bot/pkg/commands"
)

var ODFTopicAttributes = Knowledge{
	Attributes: commands.Attributes{
		MessageOfInterest: func(args []string, attribute commands.Attributes) bool {
			argMap := normalizeTokens(args)
			return tokensPresentOR(argMap, "odf")
		},
	},
	MarkdownPrompt: `ODF typically falls outside the expertise of this channel.  You might check out the ODF documentation or reach out in:
	- #forum-acm-storage.
	- #forum-rhel-coreos.
	
	Here are some resources that may help:`,
	URLS: []string{
		"<https://access.redhat.com/documentation/en-us/red_hat_openshift_data_foundation/4.14|Product Documentation for Red Hat OpenShift Data Foundation 4.14>",
	},
}
