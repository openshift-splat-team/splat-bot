package controllers

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"text/tabwriter"

	v1 "github.com/openshift-splat-team/vsphere-capacity-manager/pkg/apis/vspherecapacitymanager.splat.io/v1"
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
		Namespace: "vsphere-infra-helpers"}, pool)
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

func GetPoolStatus() (string, error) {
	var resultsBuilder strings.Builder
	poolsMu.Lock()
	defer poolsMu.Unlock()

	tbwrite := tabwriter.NewWriter(&resultsBuilder, 0, 0, 0, ' ', tabwriter.Debug)

	_, err := fmt.Fprint(tbwrite, "```\n")
	if err != nil {
		return "", err
	}
	_, err = fmt.Fprint(tbwrite, "Pool\tAvail CPUs\tAvail Mem(GB)\tSchedulable\t\n")
	if err != nil {
		return "", err
	}

	for k, v := range pools {
		_, err := fmt.Fprintf(tbwrite, "%s\t%d\t%d\t%t\t\n", k, v.Status.VCpusAvailable, v.Status.MemoryAvailable, !v.Spec.NoSchedule)
		if err != nil {
			return "", err
		}
	}
	_, err = fmt.Fprint(tbwrite, "\n```")
	if err != nil {
		return "", err
	}

	err = tbwrite.Flush()
	if err != nil {
		return "", err
	}
	return resultsBuilder.String(), nil
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
