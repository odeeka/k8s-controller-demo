package controller

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// systemNamespaces lists namespaces that are managed by Kubernetes itself.
// We never touch these.
var systemNamespaces = map[string]bool{
	"kube-system":       true,
	"kube-public":       true,
	"kube-node-lease":   true,
	"local-path-storage": true, // created by kind
}

// NamespaceLabelReconciler watches Namespace objects and ensures they carry a
// "team" label. This reconciler has no custom resource — it works directly on
// the built-in Namespace type.
type NamespaceLabelReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *NamespaceLabelReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// ── 1. Fetch the Namespace ─────────────────────────────────────────────
	// Namespaces are cluster-scoped: req.Name is the namespace name,
	// req.Namespace is always empty.
	var ns corev1.Namespace
	if err := r.Get(ctx, req.NamespacedName, &ns); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// ── 2. Skip system namespaces ──────────────────────────────────────────
	if systemNamespaces[ns.Name] {
		return ctrl.Result{}, nil
	}

	// ── 3. Check for the 'team' label ──────────────────────────────────────
	if _, hasLabel := ns.Labels["team"]; hasLabel {
		// Label already present — nothing to do.
		// Returning early here is the idempotency check.
		log.V(1).Info("Namespace already has 'team' label, skipping", "namespace", ns.Name)
		return ctrl.Result{}, nil
	}

	log.Info("Namespace is missing 'team' label, patching it", "namespace", ns.Name)

	// ── 4. Patch — add the label ───────────────────────────────────────────
	// We use MergeFrom to create a JSON Merge Patch.
	// MergeFrom takes a snapshot of the *current* state. The patch will contain
	// only the diff between the snapshot and the modified object.
	//
	// Why Patch and not Update?
	//   Update requires the full object and can conflict with concurrent writes.
	//   Patch is a targeted, atomic operation that only changes what we specify.
	base := ns.DeepCopy() // snapshot before modification
	if ns.Labels == nil {
		ns.Labels = make(map[string]string)
	}
	ns.Labels["team"] = "unassigned"

	if err := r.Patch(ctx, &ns, client.MergeFrom(base)); err != nil {
		log.Error(err, "unable to patch namespace", "namespace", ns.Name)
		return ctrl.Result{}, err
	}

	log.Info("Added 'team=unassigned' label to namespace", "namespace", ns.Name)
	return ctrl.Result{}, nil
}

// SetupWithManager registers this reconciler to watch Namespace objects.
// Note: no CRD is needed — Namespace is a standard Kubernetes type.
func (r *NamespaceLabelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Namespace{}).
		Complete(r)
}
