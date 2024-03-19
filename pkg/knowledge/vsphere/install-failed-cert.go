package vsphere

import (
	"github.com/openshift-splat-team/splat-bot/data"
	"github.com/openshift-splat-team/splat-bot/pkg/util"
)

var InstallationX509Attributes = data.Knowledge{
	Attributes: data.Attributes{
		MessageOfInterest: func(args []string, attribute data.Attributes) bool {
			argMap := util.NormalizeTokens(args)
			return util.TokensPresentOR(argMap, "vsphere", "vmware") &&
				util.TokensPresentOR(argMap, "x509") &&
				util.TokensPresentOR(argMap, "install", "installation", "ipi")
		},
	},
	MarkdownPrompt: `vCenter by default signs its own certificates. Since the vCenter certificate is not signed by a trusted CA, 
	openshift-install is unable to verify the certificate and fails. To address this, download the vCenter server root certificates[1].
	To provide the root certificate to openshift-install, provide the full path to the CA certificate in the SSL_CERT_FILE environment variable.

	For example: SSL_CERT_FILE=/path/to/vcenter-ca.crt openshift-install create cluster ...

	Here are some resources that may help:`,
	URLS: []string{
		"<https://kb.vmware.com/s/article/2108294|[1] How to download and install vCenter Server root certificates to avoid Web Browser certificate warnings>",
	},
}
