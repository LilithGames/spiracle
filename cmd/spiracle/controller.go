package main

import (
	"context"
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	// _ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	manager "sigs.k8s.io/controller-runtime/pkg/manager"

	v1 "github.com/LilithGames/spiracle/api/v1"
	"github.com/LilithGames/spiracle/controllers"
	"github.com/LilithGames/spiracle/repos"
	"github.com/LilithGames/spiracle/config"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(v1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func controller(ctx context.Context, conf *config.Config) manager.Manager {
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&zap.Options{})))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     conf.Controller.MetricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: conf.Controller.ProbeAddr,
		LeaderElection:         conf.Controller.LeaderElection.Enable,
		LeaderElectionID:       conf.Controller.LeaderElection.Id,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}
	reconciler := &controllers.RoomIngressReconciler{
		Client:    mgr.GetClient(),
		Scheme:    mgr.GetScheme(),
		Log:       ctrl.Log,
		TokenRepo: repos.NewTsTokenRepo(),
	}
	if err := reconciler.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "RoomIngress")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder
	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}
	setupLog.Info("starting manager")
	mgr.GetFieldIndexer().IndexField(ctx, &v1.RoomIngress{}, "indexToken", repos.BuildIndexToken)
	return mgr
}
