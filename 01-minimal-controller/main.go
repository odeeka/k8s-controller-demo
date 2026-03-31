// Step 01: Minimal Controller — entry point
//
// This file sets up the controller Manager and starts it.
// The Manager is the top-level object provided by controller-runtime.
// It is responsible for:
//   - Connecting to the Kubernetes API server
//   - Running our reconciler(s) in background goroutines
//   - Caching API responses so we don't hammer the API
//   - Handling OS signals (Ctrl+C → graceful shutdown)
package main

import (
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	learnv1 "github.com/example/k8s-controller-demo/01-minimal-controller/api/v1"
	"github.com/example/k8s-controller-demo/01-minimal-controller/internal/controller"
)

// scheme is the registry that maps Go types ↔ Kubernetes API group/version/kind.
// We register:
//   1. Standard Kubernetes types (Pod, Deployment, ConfigMap, …)
//   2. Our custom HelloWorld type
var scheme = runtime.NewScheme()

func init() {
	// clientgoscheme.AddToScheme registers all built-in Kubernetes types.
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	// learnv1.AddToScheme registers HelloWorld and HelloWorldList.
	utilruntime.Must(learnv1.AddToScheme(scheme))
}

func main() {
	// UseDevMode(true) produces human-readable logs instead of JSON.
	// Switch to UseDevMode(false) for production JSON logs.
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	// Create the Manager.
	// ctrl.GetConfigOrDie() reads the cluster config from:
	//   1. KUBECONFIG environment variable, or
	//   2. ~/.kube/config
	// It calls os.Exit(1) if neither is available.
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
		// Disable metrics and health probe servers — they bind to ports that
		// conflict if you run multiple steps at once. Not needed for learning.
		Metrics:                metricsserver.Options{BindAddress: "0"},
		HealthProbeBindAddress: "0",
	})
	if err != nil {
		ctrl.Log.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// Register our HelloWorld reconciler with the manager.
	// SetupWithManager configures the watch so the manager knows to
	// call Reconcile whenever a HelloWorld resource changes.
	if err := (&controller.HelloWorldReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		ctrl.Log.Error(err, "unable to create controller", "controller", "HelloWorld")
		os.Exit(1)
	}

	ctrl.Log.Info("Starting HelloWorld controller. Press Ctrl+C to stop.")

	// Start is blocking. It runs all registered controllers until it receives
	// a termination signal (SIGTERM or SIGINT / Ctrl+C).
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		ctrl.Log.Error(err, "problem running manager")
		os.Exit(1)
	}
}
