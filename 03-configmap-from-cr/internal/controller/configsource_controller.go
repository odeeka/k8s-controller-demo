package controller

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	learnv1 "github.com/example/k8s-controller-demo/03-configmap-from-cr/api/v1"
)

// ConfigSourceReconciler reconciles ConfigSource objects.
type ConfigSourceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *ConfigSourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// ── 1. Fetch the ConfigSource ──────────────────────────────────────────
	var cs learnv1.ConfigSource
	if err := r.Get(ctx, req.NamespacedName, &cs); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// ── 2. Determine the target ConfigMap name ─────────────────────────────
	cmName := cs.Spec.ConfigMapName
	if cmName == "" {
		// Default: same name as the ConfigSource
		cmName = cs.Name
	}

	log.Info("Reconciling ConfigSource", "configmap", cmName)

	// ── 3. Define the desired ConfigMap ────────────────────────────────────
	// We build the ConfigMap struct we *want* to exist.
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cmName,
			Namespace: cs.Namespace,
		},
	}

	// ── 4. CreateOrUpdate ──────────────────────────────────────────────────
	// controllerutil.CreateOrUpdate does the idempotent create-or-update for us:
	//   - If the ConfigMap doesn't exist: calls r.Create
	//   - If it already exists: calls the mutate func, then r.Update
	//
	// The mutate func receives the *existing* object (or an empty one for
	// creates) and applies our desired state to it. We only set the Data field,
	// leaving all other fields (resourceVersion, labels, etc.) untouched.
	op, err := controllerutil.CreateOrUpdate(ctx, r.Client, cm, func() error {
		cm.Data = cs.Spec.Data
		return nil
	})
	if err != nil {
		log.Error(err, "unable to create or update ConfigMap")
		return ctrl.Result{}, err
	}
	log.Info("ConfigMap reconciled", "operation", op, "name", cmName)

	// ── 5. Update status ───────────────────────────────────────────────────
	cs.Status.Phase = "Ready"
	cs.Status.ManagedConfigMap = cmName
	if err := r.Status().Update(ctx, &cs); err != nil {
		log.Error(err, "unable to update ConfigSource status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *ConfigSourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&learnv1.ConfigSource{}).
		Complete(r)
}
