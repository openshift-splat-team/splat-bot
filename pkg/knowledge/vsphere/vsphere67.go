package vsphere

import (
	"github.com/openshift-splat-team/splat-bot/data"
	"github.com/openshift-splat-team/splat-bot/pkg/util"
)

var VSphere67TopicAttributes = data.Knowledge{
	Attributes: data.Attributes{
		MessageOfInterest: func(args []string, attribute data.Attributes) bool {
			argMap := util.NormalizeTokens(args)
			return util.TokensPresentOR(argMap, "vsphere", "vmware", "vcenter") && util.TokensPresentAND(argMap, "6.7")
		},
	},
	MarkdownPrompt: `vSphere 6 - 6.7 is end of life. Unless a customer has a VMware extended support agreement and a RH support exception OpenShift running on that version is not supported.
	Please provide:
	- Support case number
	- OpenShift and vSphere versions
	- Appropriate logs - this consists of either a must-gather or install-gather`,
	URLS: []string{
		"<https://docs.openshift.com/container-platform/latest/support/gathering-cluster-data.html#about-must-gather_gathering-cluster-data|gathering a must-gather>",
		"<https://docs.openshift.com/container-platform/latest/support/troubleshooting/troubleshooting-installations.html#installation-bootstrap-gather_troubleshooting-installations|gathering an install-gather>",
	},
}
