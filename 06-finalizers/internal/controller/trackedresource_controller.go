package controller

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	learnv1 "github.com/example/k8s-controller-demo/06-finalizers/api/v1"
)

// cleanupFinalizer is the finalizer string we add to every TrackedResource.
// Convention: use <group>/<descriptive-name>.
const cleanupFinalizer = "learn.example.com/cleanup"

// TrackedResourceReconciler reconciles TrackedResource objects.
type TrackedResourceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *TrackedResourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// ── 1. Fetch the TrackedResource ───────────────────────────────────────
	var tr learnv1.TrackedResource
	if err := r.Get(ctx, req.NamespacedName, &tr); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// ── 2. Check if the resource is being deleted ──────────────────────────
	// When a resource has a DeletionTimestamp, it means `kubectl delete` was
	// called but Kubernetes is waiting for all finalizers to be removed.
	if !tr.DeletionTimestamp.IsZero() {
		return r.handleDeletion(ctx, &tr)
	}

	// ── 3. Normal reconciliation ───────────────────────────────────────────
	return r.handleNormal(ctx, &tr)
}

// handleNormal is called when the resource is alive (not being deleted).
func (r *TrackedResourceReconciler) handleNormal(ctx context.Context, tr *learnv1.TrackedResource) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Ensure the finalizer is registered.
	// If it's not present, add it and return — the update will trigger another
	// reconcile so we continue from a clean state.
	if !controllerutil.ContainsFinalizer(tr, cleanupFinalizer) {
		log.Info("Adding finalizer to TrackedResource", "name", tr.Name)
		controllerutil.AddFinalizer(tr, cleanupFinalizer)
		if err := r.Update(ctx, tr); err != nil {
			return ctrl.Result{}, err
		}
		// Return here — the Update will trigger a new reconcile.
		return ctrl.Result{}, nil
	}

	// Normal work: update status to show the resource is active.
	log.Info("Reconciling TrackedResource", "message", tr.Spec.Message)

	tr.Status.Phase = "Active"
	tr.Status.Message = fmt.Sprintf("Tracking: %s", tr.Spec.Message)
	if err := r.Status().Update(ctx, tr); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// handleDeletion is called when DeletionTimestamp is set (kubectl delete was called).
// It runs cleanup and then removes the finalizer so Kubernetes can finish deleting.
func (r *TrackedResourceReconciler) handleDeletion(ctx context.Context, tr *learnv1.TrackedResource) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	if !controllerutil.ContainsFinalizer(tr, cleanupFinalizer) {
		// Finalizer already removed — nothing to do.
		return ctrl.Result{}, nil
	}

	log.Info("DeletionTimestamp is set — running cleanup", "name", tr.Name, "message", tr.Spec.Message)

	// ── Perform your cleanup here ──────────────────────────────────────────
	// In a real controller this might be:
	//   - Deleting an external database record
	//   - Releasing a DNS entry
	//   - Deregistering from an external service
	//   - Waiting for a cloud resource to be deleted
	//
	// For this example we just log. If cleanup fails, return an error and
	// the framework will requeue with backoff.
	if err := r.doCleanup(ctx, tr); err != nil {
		log.Error(err, "cleanup failed, will retry")
		return ctrl.Result{}, err
	}

	// ── Remove the finalizer ───────────────────────────────────────────────
	// Once the finalizer is removed, Kubernetes will complete the deletion.
	// Do this LAST — after you are certain cleanup succeeded.
	log.Info("Cleanup complete, removing finalizer", "name", tr.Name)
	controllerutil.RemoveFinalizer(tr, cleanupFinalizer)
	if err := r.Update(ctx, tr); err != nil {
		return ctrl.Result{}, err
	}

	log.Info("Finalizer removed; Kubernetes will now delete the resource", "name", tr.Name)
	return ctrl.Result{}, nil
}

// doCleanup simulates performing external cleanup.
// Replace the body with your actual cleanup logic.
func (r *TrackedResourceReconciler) doCleanup(ctx context.Context, tr *learnv1.TrackedResource) error {
	log := log.FromContext(ctx)
	log.Info("Simulating external cleanup...",
		"resource", tr.Name,
		"message", tr.Spec.Message,
	)
	// Simulate success. In a real scenario, errors here cause requeue + retry.
	return nil
}

func (r *TrackedResourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&learnv1.TrackedResource{}).
		Complete(r)
}
