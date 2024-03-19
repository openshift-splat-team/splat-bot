package vsphere

import (
	"github.com/openshift-splat-team/splat-bot/data"
	"github.com/openshift-splat-team/splat-bot/pkg/util"
)

var ODFTopicAttributes = data.Knowledge{
	Attributes: data.Attributes{
		MessageOfInterest: func(args []string, attribute data.Attributes) bool {
			argMap := util.NormalizeTokens(args)
			return util.TokensPresentOR(argMap, "odf")
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
