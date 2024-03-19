package vsphere

import (
	"strings"
	"testing"
)

const (
	x509InstallFailed = `Hi Team.
	I am trying to install OCP 4.15 using IPI on a fresh "VMware Cloud Public Cloud Open Environment" demo item. I am facing a TLS issue:
	[Downloads]$ ./openshift-install create cluster --dir ocpinstall
	INFO Consuming Install Config from target directory 
	INFO Creating infrastructure resources...         
	ERROR                                              
	ERROR Error: failed to upload: Post "https://x.x.x.x/nfc/xxxxxxxx-xxxxxx-xxxxxx/disk-0.vmdk": tls: failed to verify certificate: x509: certificate is valid for 10.0.0.1, not 20.0.0.1 
	ERROR                                              
	ERROR   with vsphereprivate_import_ova.import["generated-failure-domain"], 
	ERROR   on main.tf line 63, in resource "vsphereprivate_import_ova" "import": 
	ERROR   63: resource "vsphereprivate_import_ova" "import" { 
	ERROR                                              
	ERROR failed to fetch Cluster: failed to generate asset "Cluster": failed to create cluster: failure applying terraform for "pre-bootstrap" stage: error applying Terraform configs: failed to apply Terraform: exit status 1 
	ERROR                                              
	ERROR Error: failed to upload: Post "https://x.x.x.x/nfc/xxxxxxxx-xxxxxx-xxxxxx/disk-0.vmdk": tls: failed to verify certificate: x509: certificate is valid for 10.0.0.1, not 20.0.0.1 
	ERROR                                              
	ERROR   with vsphereprivate_import_ova.import["generated-failure-domain"], 
	ERROR   on main.tf line 63, in resource "vsphereprivate_import_ova" "import": 
	ERROR   63: resource "vsphereprivate_import_ova" "import" { 
	ERROR`
)

func TestInstallfailed(t *testing.T) {
	tokens := strings.Split(x509InstallFailed, " ")

	if !InstallationX509Attributes.MessageOfInterest(tokens, InstallationX509Attributes.Attributes) {
		t.Errorf("expected to match message of interest")
	}
}
