package controller

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	learnv1 "github.com/example/k8s-controller-demo/04-deployment-manager/api/v1"
)

// AppDeploymentReconciler reconciles AppDeployment objects.
type AppDeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *AppDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// ── 1. Fetch the AppDeployment ─────────────────────────────────────────
	var app learnv1.AppDeployment
	if err := r.Get(ctx, req.NamespacedName, &app); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("Reconciling AppDeployment", "image", app.Spec.Image, "replicas", app.Spec.Replicas)

	// ── 2. Prepare the Deployment shell ───────────────────────────────────
	// We only set Name + Namespace here — this is the "lookup key" that
	// CreateOrUpdate uses to GET the existing Deployment (or create a new one).
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
		},
	}

	// ── 3. CreateOrUpdate the Deployment ───────────────────────────────────
	// IMPORTANT: All desired-state mutations go INSIDE the mutate function.
	//
	// Why? CreateOrUpdate works in two steps:
	//   1. GET the existing object from the API (or start with the shell above)
	//   2. Call this function so you can set desired fields on it
	//   3. Create or Update the result
	//
	// If you set fields *before* CreateOrUpdate, those values are overwritten
	// by the GET in step 1 — so your changes are lost on the update path.
	// This includes owner references, labels, and spec fields.
	result, err := controllerutil.CreateOrUpdate(ctx, r.Client, deploy, func() error {
		// ── Owner reference ──────────────────────────────────────────────
		// Links `deploy` to `app`. When `app` is deleted, Kubernetes will
		// automatically garbage-collect `deploy`.
		// Must be set inside the mutate func so it survives the update path.
		if err := controllerutil.SetControllerReference(&app, deploy, r.Scheme); err != nil {
			return err
		}

		// ── Desired spec ──────────────────────────────────────────────────
		replicas := app.Spec.Replicas
		labels := map[string]string{
			"app":                          app.Name,
			"app.kubernetes.io/managed-by": "appdeployment-controller",
		}

		deploy.Spec.Replicas = &replicas
		deploy.Spec.Selector = &metav1.LabelSelector{MatchLabels: labels}
		deploy.Spec.Template = corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{Labels: labels},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{{
					Name:  app.Name,
					Image: app.Spec.Image,
					Ports: []corev1.ContainerPort{{ContainerPort: app.Spec.Port}},
				}},
			},
		}
		return nil
	})
	if err != nil {
		log.Error(err, "unable to create or update Deployment")
		return ctrl.Result{}, err
	}
	log.Info("Deployment reconciled", "operation", result)

	// ── 5. Read the Deployment back to get its current status ──────────────
	// The Deployment status is filled in by the Deployment controller.
	// We need to re-fetch it to get the latest availableReplicas.
	if err := r.Get(ctx, req.NamespacedName, deploy); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// ── 6. Update AppDeployment status ─────────────────────────────────────
	app.Status.AvailableReplicas = deploy.Status.AvailableReplicas
	if deploy.Status.AvailableReplicas >= app.Spec.Replicas {
		app.Status.Phase = "Available"
	} else {
		app.Status.Phase = "Progressing"
	}

	if err := r.Status().Update(ctx, &app); err != nil {
		log.Error(err, "unable to update AppDeployment status")
		return ctrl.Result{}, err
	}

	log.Info("Status updated",
		"phase", app.Status.Phase,
		"availableReplicas", app.Status.AvailableReplicas,
	)
	return ctrl.Result{}, nil
}

// SetupWithManager registers this reconciler.
// Note the Owns(&appsv1.Deployment{}) — this triggers a reconcile for the
// parent AppDeployment whenever an owned Deployment changes (e.g., replicas
// become available). Without this, status.availableReplicas would never update.
func (r *AppDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&learnv1.AppDeployment{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
