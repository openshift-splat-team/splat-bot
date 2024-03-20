package vsphere

import (
	"strings"
	"testing"
)

const (
	migrationQuestion = `Hi Team, One of my xyz customers, has a requirement to migrate nn+ clusters from vsphere 6 to vsphere 7 so they can unblock themselves, allowing the upgrading of their OCP4 clusters from OCP4.12 to OCP4.13 etc
	They are not in a position to upgrade their vsphere 6 clusters, but instead, have built new vsphere 7 clusters and have tested migrating one of their Pilot clusters to the new vcentre which appears to have succeeded without any issues after they updated the vsphere-creds secret and also the cloud-provider-config configmap.
	All nodes are on-line and are reporting as passed in the latest MG.`
)

func TestMigrationMesssage(t *testing.T) {
	tokens := strings.Split(migrationQuestion, " ")

	if !MigrationTopicAttributes.MessageOfInterest(tokens, InstallationX509Attributes.Attributes) {
		t.Errorf("expected to match message of interest")
	}
}
