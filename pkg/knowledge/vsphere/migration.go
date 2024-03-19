package vsphere

import (
	"github.com/openshift-splat-team/splat-bot/data"
	"github.com/openshift-splat-team/splat-bot/pkg/util"
	knowledge "github.com/openshift-splat-team/splat-bot/pkg/util"
)

var MigrationTopicAttributes = data.Knowledge{
	Attributes: data.Attributes{
		MessageOfInterest: func(args []string, attribute data.Attributes) bool {
			argMap := util.NormalizeTokens(args)
			if knowledge.TokensPresentOR(argMap, "migration", "vmotion") &&
				knowledge.TokensPresentOR(argMap, "vsphere", "vmware") {
				return true
			}
			return false
		},
	},
	MarkdownPrompt: `Migration can be a complex topic as it relates to vSphere and OpenShift.  
	The TL;DR:
	- Storage vMotion isn't supported, however, compute vMotion is
	- Migration of VMs between vCenters isn't supported
	
	Here are some resources that may help:`,
	URLS: []string{
		"<https://access.redhat.com/solutions/6509731|Migrating directly between vCenters is not supported>",
		"<https://docs.openshift.com/container-platform/4.15/migration_toolkit_for_containers/about-mtc.html|Migration Toolkit for Containers>",
		"<https://access.redhat.com/articles/6718991|Migrating Virtual Machines with vMotion and the Machine API>",
	},
}
