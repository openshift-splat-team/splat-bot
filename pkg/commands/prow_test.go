package commands

import "testing"

func TestProw(t *testing.T) {
	startProwRetrievalTimers()
	createProwGraph("vsphere")
	
}