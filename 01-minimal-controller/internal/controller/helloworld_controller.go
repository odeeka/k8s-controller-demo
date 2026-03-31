package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	learnv1 "github.com/example/k8s-controller-demo/01-minimal-controller/api/v1"
)

// HelloWorldReconciler is the controller for HelloWorld resources.
//
// A Reconciler is just a struct with a Reconcile method.  It usually holds:
//   - client.Client  — to talk to the Kubernetes API
//   - *runtime.Scheme — to create and decode objects
//
// Additional fields (like loggers, configuration, or external clients) can
// be added here as your controller grows.
type HelloWorldReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Reconcile is called by controller-runtime whenever a HelloWorld resource
// is created, updated, or deleted, and also at startup for all existing ones.
//
// The framework guarantees:
//   - At-least-once delivery: if a change is missed, it will be retried.
//   - Serialised per-object: two Reconcile calls for the same resource
//     are never running concurrently.
func (r *HelloWorldReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// log.FromContext returns a structured logger pre-populated with the
	// reconciler name and the resource name/namespace.
	log := log.FromContext(ctx)

	// ── Step 1: Fetch the resource ─────────────────────────────────────────
	// req.NamespacedName is just the name+namespace of the changed resource.
	// We need to call r.Get to load the full object from the cache.
	var hw learnv1.HelloWorld
	if err := r.Get(ctx, req.NamespacedName, &hw); err != nil {
		// client.IgnoreNotFound converts a "not found" error into nil.
		// This happens when the resource was deleted before we got to reconcile it.
		// Returning nil tells the framework: "nothing to do, don't requeue".
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// ── Step 2: Do something ───────────────────────────────────────────────
	// In this minimal example we just log. In later steps, we'll create
	// Kubernetes resources or update the status here.
	log.Info("Reconciling HelloWorld",
		"spec.name", hw.Spec.Name,
		"resource", req.NamespacedName,
	)

	// ── Step 3: Return ─────────────────────────────────────────────────────
	// ctrl.Result{} means "success, don't requeue — wait for the next event".
	return ctrl.Result{}, nil
}

// SetupWithManager registers this reconciler with the Manager and declares
// which resources it wants to watch.
//
// ctrl.NewControllerManagedBy(mgr).For(&learnv1.HelloWorld{}) means:
//   "call my Reconcile method whenever a HelloWorld resource changes."
func (r *HelloWorldReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&learnv1.HelloWorld{}).
		Complete(r)
}
