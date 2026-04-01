package controller

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	learnv1 "github.com/example/k8s-controller-demo/02-status-updates/api/v1"
)

// GreeterReconciler reconciles Greeter objects.
type GreeterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *GreeterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// ── 1. Fetch the Greeter ───────────────────────────────────────────────
	var gr learnv1.Greeter
	if err := r.Get(ctx, req.NamespacedName, &gr); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("Reconciling Greeter",
		"greeting", gr.Spec.Greeting,
		"target", gr.Spec.TargetName,
	)

	// ── 2. Compute desired state ───────────────────────────────────────────
	desiredMessage := fmt.Sprintf("%s, %s!", gr.Spec.Greeting, gr.Spec.TargetName)

	// ── 3. Idempotency check ───────────────────────────────────────────────
	// Only write to the API if something actually changed.
	// Without this check, every reconcile would issue an Update call, which
	// would trigger another reconcile event — creating a tight loop.
	if gr.Status.Phase == "Ready" && gr.Status.Message == desiredMessage {
		log.Info("Status already up to date, nothing to do")
		return ctrl.Result{}, nil
	}

	// ── 4. Update status ───────────────────────────────────────────────────
	// IMPORTANT: Use r.Status().Update(), NOT r.Update().
	//
	// When the status subresource is enabled (see the CRD yaml), the API server
	// splits the resource into two independent endpoints:
	//   PUT /apis/learn.example.com/v1/namespaces/*/greeters/<name>        → updates spec
	//   PUT /apis/learn.example.com/v1/namespaces/*/greeters/<name>/status → updates status
	//
	// r.Update()        → writes to the spec endpoint (status changes ignored)
	// r.Status().Update() → writes to the status endpoint (spec changes ignored)
	//
	// This separation prevents a user's `kubectl edit` from accidentally
	// overwriting controller-managed status fields.
	gr.Status.Phase = "Ready"
	gr.Status.Message = desiredMessage
	gr.Status.LastUpdatedTime = metav1.Now()

	if err := r.Status().Update(ctx, &gr); err != nil {
		log.Error(err, "unable to update Greeter status")
		return ctrl.Result{}, err
	}

	log.Info("Greeter status updated", "message", desiredMessage)
	return ctrl.Result{}, nil
}

func (r *GreeterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&learnv1.Greeter{}).
		Complete(r)
}
