# Step 04 — Deployment Manager

**New concepts:** Owner references · Watching owned resources · Syncing child resource state

---

## What This Step Does

We introduce an `AppDeployment` custom resource that describes a simple application. The controller:

1. Creates a `Deployment` based on `spec.image` and `spec.replicas`
2. Sets an **owner reference** from the `Deployment` to the `AppDeployment`
3. Watches both `AppDeployment` and its owned `Deployment` objects
4. Updates `status.availableReplicas` by reading the Deployment's own status

---

## Key Ideas

### Owner references

An **owner reference** links a child resource (the `Deployment`) to its parent (the `AppDeployment`). When the parent is deleted, Kubernetes automatically garbage-collects all children. This prevents resource leaks.

```go
// Set an owner reference: when `app` is deleted, `deploy` is garbage-collected.
controllerutil.SetControllerReference(&app, &deploy, r.Scheme)
```

The reference contains:
- The parent's UID (to uniquely identify it)
- `controller: true` (so only one controller "owns" the resource)
- `blockOwnerDeletion: true` (lets the parent block deletion until children are gone)

### Watching owned resources

```go
ctrl.NewControllerManagedBy(mgr).
    For(&learnv1.AppDeployment{}).
    Owns(&appsv1.Deployment{}).   // ← triggers reconcile when a Deployment changes
    Complete(r)
```

`Owns` tells the manager: if a `Deployment` that is owned by an `AppDeployment` changes (e.g., Pods become Ready), trigger a reconcile for that parent `AppDeployment`. This is how we can update `status.availableReplicas` without polling.

### Desired-state comparison

Rather than tracking what we created before, we simply re-compute the desired state on every reconcile and compare it to what's currently in the cluster. If they differ, we update. This is the foundation of **declarative** control.

---

## Project Layout

```
04-deployment-manager/
├── main.go
├── api/v1/
│   ├── groupversion_info.go
│   ├── appdeployment_types.go   # spec.image, spec.replicas, status
│   └── zz_generated.deepcopy.go
├── internal/controller/
│   └── appdeployment_controller.go  # CreateOrUpdate + owner ref + status sync
└── config/
    ├── crd/learn.example.com_appdeployments.yaml
    └── samples/appdeployment_sample.yaml
```

---

## How to Run

```bash
kubectl apply -f config/crd/
kubectl apply -f config/samples/
go run .
```

Watch the Deployment get created:
```bash
kubectl get appdeployment my-app
# NAME     REPLICAS   AVAILABLE   PHASE
# my-app   2          2           Available

kubectl get deployment my-app
# Shows the managed Deployment
```

### Test owner reference cleanup
```bash
kubectl delete appdeployment my-app
kubectl get deployment my-app
# Error: deployment "my-app" not found  — auto-cleaned up!
```

### Test spec changes
```bash
kubectl patch appdeployment my-app --type=merge -p '{"spec":{"replicas":3}}'
kubectl get deployment my-app
# DESIRED 3 — Deployment was updated automatically
```

---

## What's Next

[Step 05 →](../05-label-enforcer/) Watching built-in resources without a CRD.
