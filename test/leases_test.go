package test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift-splat-team/splat-bot/pkg/controllers"
	v1 "github.com/openshift-splat-team/vsphere-capacity-manager/pkg/apis/vspherecapacitymanager.splat.io/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	k8sctrl "sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

var _ = Describe("Lease management", func() {
	timeout := 10 * time.Second

	It("should be able to get pool name", func() {
		names, err := controllers.GetPoolNames(ctx)
		Expect(err).To(BeNil(), "No error should occur when getting the list of pool names")
		Expect(names).To(HaveLen(4), "There should be 4 pools")
	})

	It("should be able to handle a simple lease", func() {
		user := "user1"
		By("acquiring it", func() {
			leases, err := controllers.AcquireLease(ctx, user, 1, 1, "", 1)
			Expect(err).To(BeNil())
			Expect(leases).NotTo(BeNil())
		})

		By("checking the lease's status", func() {
			// Currently this returns a complex string.  We can have an expected output to compare against, or just
			// be happy if its not len() == 0
			status, err := controllers.GetLeaseStatus(user)
			Expect(err).To(BeNil())
			Expect(status).NotTo(HaveLen(0))
		})

		By("checking the lease exists", func() {
			Eventually(func() int {
				return len(getLeases(mgrClient, user, true))
			}, timeout).Should(Equal(1))
		})

		By("releasing it", func() {
			err := controllers.RemoveLease(ctx, user)
			Expect(err).To(BeNil())
		})

		By("verifying lease is gone", func() {
			Eventually(func() int {
				return len(getLeases(mgrClient, user, true))
			}, timeout).Should(Equal(0))
		})
	})

	It("should be able to create a lease with numerous additional networks", func() {
		user := "user2"
		By("acquiring it", func() {
			leases, err := controllers.AcquireLease(ctx, user, 1, 1, "", 4)
			Expect(err).To(BeNil())
			Expect(leases).NotTo(BeNil())
		})

		By("checking the lease's status", func() {
			// Currently this returns a complex string.  We can have an expected output to compare against, or just
			// be happy if its not len() == 0
			status, err := controllers.GetLeaseStatus(user)
			Expect(err).To(BeNil())
			Expect(status).NotTo(HaveLen(0))
		})

		// For this test, we need to actually check the number of leases
		By("verifying real count of leases", func() {
			Eventually(func() int {
				return len(getLeases(mgrClient, user, true))
			}, timeout).Should(Equal(4))
		})

		By("releasing it", func() {
			err := controllers.RemoveLease(ctx, user)
			Expect(err).To(BeNil())
		})

		By("verifying lease is gone", func() {
			Eventually(func() int {
				return len(getLeases(mgrClient, user, true))
			}, timeout).Should(Equal(0))
		})
	})
})

func getLeases(mgrClient k8sctrl.Client, user string, includeNetworkOnly bool) []v1.Lease {
	//fmt.Printf("Getting leases for user %v\n", user)
	leases := &v1.LeaseList{}

	var labelSelector labels.Selector
	if !includeNetworkOnly {
		userLeaseReq, err := labels.NewRequirement(controllers.SplatBotLeaseOwner, selection.Equals, []string{user})
		Expect(err).To(BeNil())
		networkLeaseReq, err := labels.NewRequirement("network-only-lease", selection.DoesNotExist, []string{})
		Expect(err).To(BeNil())

		labelSelector = labels.NewSelector()
		labelSelector.Add(*userLeaseReq, *networkLeaseReq)
	} else {
		labelSelector = labels.SelectorFromSet(labels.Set{controllers.SplatBotLeaseOwner: user})
	}

	listOptions := &k8sctrl.ListOptions{
		LabelSelector: labelSelector,
		Namespace:     controllers.VcmNamespace,
	}

	Expect(mgrClient.List(ctx, leases, listOptions)).To(Succeed())
	return leases.Items
}
