package controllers

import (
	"context"
	"fmt"
	"strings"
	"sync"

	v1 "github.com/openshift-splat-team/vsphere-capacity-manager/pkg/apis/vspherecapacitymanager.splat.io/v1"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	poolsMu   sync.Mutex
	pools     = make(map[string]*v1.Pool)
	k8sclient client.Client
)

func SetPoolSchedulable(ctx context.Context, name string, schedulable bool) error {
	pool := &v1.Pool{}
	name = strings.ReplaceAll(name, "<http://vcenter.ci|vcenter.ci>", "vcenter.ci")
	err := k8sclient.Get(ctx, types.NamespacedName{
		Name:      name,
		Namespace: VcmNamespace}, pool)
	if err != nil {
		return fmt.Errorf("could not get pool %s: %v", name, err)
	}
	pool.Spec.NoSchedule = !schedulable
	err = k8sclient.Update(ctx, pool)
	if err != nil {
		return fmt.Errorf("could not update pool %s: %v", name, err)
	}
	return nil
}

func GetPoolStatus() (slack.MsgOption, error) {
	poolsMu.Lock()
	defer poolsMu.Unlock()

	var rtBlocks []slack.Block

	for idx, pool := range pools {
		var rtElems []slack.RichTextElement

		availCPU := float64(100) * float64(pool.Status.VCpusAvailable) / float64(pool.Spec.VCpus)
		availMemory := float64(100) * float64(pool.Status.MemoryAvailable) / float64(pool.Spec.Memory)
		status := fmt.Sprintf("\tCPU: %.0f%%, Memory: %.0f%%", availCPU, availMemory)
		color := "large_green_circle"
		if pool.Spec.NoSchedule {
			status = fmt.Sprintf("\t!! Cordoned !! %s", status)
			color = "black_circle"
		} else {
			if availCPU < 20 {
				color = "red_circle"
			} else if availMemory < 20 {
				color = "red_circle"
			} else if availCPU < 50 {
				color = "large_yellow_circle"
			} else if availMemory < 50 {
				color = "large_yellow_circle"
			}
		}

		rtElems = append(rtElems, slack.NewRichTextSection([]slack.RichTextSectionElement{
			slack.NewRichTextSectionEmojiElement(color, 2, nil),
			slack.NewRichTextSectionTextElement(fmt.Sprintf(" %s\n", pool.Name), &slack.RichTextSectionTextStyle{Bold: true}),
			slack.NewRichTextSectionTextElement(status, &slack.RichTextSectionTextStyle{}),
		}...))

		rtBlocks = append(rtBlocks, slack.NewRichTextBlock(fmt.Sprintf("pool-status-%s", idx), rtElems...))
		/*rtBlocks = append(rtBlocks, slack.NewSectionBlock(slack.NewTextBlockObject(slack.PlainTextType, ".", false, false), nil, slack.NewAccessory(
		slack.NewButtonBlockElement(fmt.Sprintf("cordon-%s", idx), "cordon", slack.NewTextBlockObject(slack.PlainTextType, "cordon", false, false)))))*/
		rtBlocks = append(rtBlocks, slack.NewDividerBlock())
	}

	blocks := []slack.Block{
		slack.NewHeaderBlock(slack.NewTextBlockObject(slack.PlainTextType, "CI Pool Status", false, false)),
		slack.NewSectionBlock(slack.NewTextBlockObject(slack.PlainTextType, "Pool status shows free CPU and memory as well as the state of the pools.\n\n", false, false), nil, nil),
		slack.NewDividerBlock(),
	}

	blocks = append(blocks, rtBlocks...)

	return slack.MsgOptionBlocks(blocks...), nil
}

type PoolReconciler struct {
	client.Client
	Scheme         *runtime.Scheme
	Recorder       record.EventRecorder
	RESTMapper     meta.RESTMapper
	UncachedClient client.Client

	// Namespace is the namespace in which the ControlPlaneMachineSet controller should operate.
	// Any ControlPlaneMachineSet not in this namespace should be ignored.
	Namespace string

	// OperatorName is the name of the ClusterOperator with which the controller should report
	// its status.
	OperatorName string

	// ReleaseVersion is the version of current cluster operator release.
	ReleaseVersion string
}

func (l *PoolReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := ctrl.NewControllerManagedBy(mgr).
		For(&v1.Pool{}).
		Complete(l); err != nil {
		return fmt.Errorf("error setting up controller: %w", err)
	}

	// Set up API helpers from the manager.
	l.Client = mgr.GetClient()
	k8sclient = mgr.GetClient()
	l.Scheme = mgr.GetScheme()
	l.Recorder = mgr.GetEventRecorderFor("pools-controller")
	l.RESTMapper = mgr.GetRESTMapper()

	return nil
}

func (l *PoolReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log.Print("Reconciling pool")
	defer log.Print("Finished reconciling pool")
	pool := &v1.Pool{}
	err := l.Client.Get(ctx, types.NamespacedName{
		Namespace: req.Namespace,
		Name:      req.Name,
	}, pool)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	poolsMu.Lock()
	pools[pool.Name] = pool
	poolsMu.Unlock()

	return ctrl.Result{}, nil
}
