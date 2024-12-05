package controllers

import (
	log "github.com/sirupsen/logrus"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"

	v1 "github.com/openshift-splat-team/vsphere-capacity-manager/pkg/apis/vspherecapacitymanager.splat.io/v1"
	"k8s.io/klog/v2/textlogger"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func init() {

	if os.Getenv("UNIT") != "" {
		log.Printf("!!! controllers are disabled for unit tests")
		return
	}
	logger := textlogger.NewLogger(textlogger.NewConfig())
	ctrl.SetLogger(logger)

	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		log.Printf("could not create manager: %v", err)
		os.Exit(1)
	}

	err = v1.AddToScheme(mgr.GetScheme())
	if err != nil {
		log.Printf("could not add types to scheme: %v", err)
		os.Exit(1)
	}

	if err := (&PoolReconciler{}).
		SetupWithManager(mgr); err != nil {
		log.Printf("unable to create controller: %v", err)
		os.Exit(1)
	}

	if err := (&LeaseReconciler{}).
		SetupWithManager(mgr); err != nil {
		log.Printf("unable to create controller: %v", err)
		os.Exit(1)
	}

	go func() {
		if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
			log.Printf("could not start manager: %v", err)
			os.Exit(1)
		}
	}()
}
