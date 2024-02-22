package events

import (
	"fmt"
	"testing"
)

var questions = []string{
	"my customer wants to make it clear whether a reverse DNS record(PTR record) for API VIP is necessary or not in vSphere IPI installation.\nIn our document, only A/AAAA or CNAME record is mentioned so I am assuming that PTR record is not necessary but I would like to double-check it here just in case. Thank you!",
	"Hello Team #forum-vmware just need a quick confirmation regarding vSphere 6.7 compatibility with OCP4.10. As per kcs: https://access.redhat.com/articles/4763741 vSphere 6.7 is not tested on OCP4.10. but as per doc[1], vSphere 6.5 or later is supported. Can someone please confirm the suppotability here?\n",
	"IHAC who is looking for the supported way to gracefully move the nodes to the new datastore?  in Disconnected/Isolated environment, they are refreshing storage Nimble hardware and datastore names will be changed.",
	"I know that SDRS (dynamic storage migration) is not supported with OpenShift. I read that this may be related to the in-tree driver, but wasn't able to find details on this. I assume this is likely only related to pod-requested storage and the node's OS is handled properly - otherwise VMWare would have a lot of problems with it.\nQ: could anyone provide some details why that is not supported?\nQ2: would this statement change if the CSI driver is in use?",
	"Do monkeys like bananas?",
}

func TestTransform(t *testing.T) {
	err := Init()
	if err != nil {
		t.Error(err)
		return
	}

	for _, question := range questions {

		fmt.Printf("question: %s\n", question)
		_, err = Response(question)
		if err != nil {
			t.Error(err)
			return
		}
	}
}
