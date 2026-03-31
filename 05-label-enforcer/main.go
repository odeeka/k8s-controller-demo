// Step 05: Label Enforcer — entry point
//
// Notice: there is no custom API type to register.
// We only need the standard clientgoscheme to work with Namespaces.
package main

import (
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/example/k8s-controller-demo/05-label-enforcer/internal/controller"
)

var scheme = runtime.NewScheme()

func init() {
	// We only need standard Kubernetes types — no custom CRD in this step.
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
}

func main() {
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
	})
	if err != nil {
		ctrl.Log.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err := (&controller.NamespaceLabelReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		ctrl.Log.Error(err, "unable to create controller", "controller", "NamespaceLabel")
		os.Exit(1)
	}

	ctrl.Log.Info("Starting Namespace label enforcer. Press Ctrl+C to stop.")

	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		ctrl.Log.Error(err, "problem running manager")
		os.Exit(1)
	}
}
