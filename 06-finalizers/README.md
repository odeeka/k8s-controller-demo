# Step 06 — Finalizers

**New concepts:** Finalizers · Deletion lifecycle · Graceful cleanup

---

## What This Step Does

We introduce a `TrackedResource` custom resource that simulates a resource needing cleanup before deletion (e.g., an external registration, a lease, a database entry).

The controller:

1. Adds a **finalizer** to every new `TrackedResource`
2. On deletion: detects the `DeletionTimestamp`, runs cleanup, then removes the finalizer
3. After the finalizer is removed, Kubernetes completes the deletion

---

## Key Ideas

### Why finalizers?

When you delete a Kubernetes resource normally, it is immediately removed from etcd. But what if your controller created something *outside* Kubernetes (an external database row, a DNS record, a cloud resource)? You need a way to run cleanup code before the resource disappears.

Finalizers are strings in `metadata.finalizers`. While any finalizer strings exist, Kubernetes will NOT delete the resource — it sets `metadata.deletionTimestamp` and waits. Your controller is responsible for doing cleanup and then removing the finalizer string.

### The deletion flow

```text
User: kubectl delete trackedresource my-resource
         │
         ▼
Kubernetes sets DeletionTimestamp, resource stays in etcd
         │
         ▼ (triggers reconcile)
Controller sees DeletionTimestamp is set
         │
         ├── Runs cleanup (delete external record, release lease, etc.)
         │
         └── Removes finalizer from metadata.finalizers
                   │
                   ▼
         Kubernetes deletes the resource from etcd ✓
```

### Reading the code

The reconciler has two distinct branches:

```go
if !resource.DeletionTimestamp.IsZero() {
    // ── DELETION PATH ──────────────────────────────────────
    // The resource is being deleted.
    // Run cleanup, then remove the finalizer.
} else {
    // ── NORMAL PATH ────────────────────────────────────────
    // The resource is alive.
    // Ensure the finalizer is registered.
    // Do normal work.
}
```

---

## Project Layout

```text
06-finalizers/
├── main.go
├── api/v1/
│   ├── groupversion_info.go
│   ├── trackedresource_types.go
│   └── zz_generated.deepcopy.go
├── internal/controller/
│   └── trackedresource_controller.go  # ← finalizer add/remove pattern
└── config/
    ├── crd/learn.example.com_trackedresources.yaml
    └── samples/trackedresource_sample.yaml
```

---

## How to Run

```bash
kubectl apply -f config/crd/
kubectl apply -f config/samples/
go run .
```

Observe the finalizer being added:

```bash
kubectl get trackedresource my-tracked -o yaml
# metadata.finalizers:
#   - learn.example.com/cleanup
```

Delete the resource and watch the logs:

```bash
kubectl delete trackedresource my-tracked
# Controller logs:
#   INFO  DeletionTimestamp set, running cleanup ...
#   INFO  Cleanup complete, removing finalizer ...
#   INFO  Finalizer removed; Kubernetes will now delete the resource
```

### What happens if you force-remove the finalizer?

```bash
# This bypasses your cleanup logic — use with caution!
kubectl patch trackedresource my-tracked \
  --type=json -p '[{"op":"remove","path":"/metadata/finalizers"}]'
```

## Run In-Cluster (Deployment)

If you want this controller to run inside Kubernetes, use the manifests under `config/rbac/` and `config/manager/`.

### 1. Build and push image

From repository root (`k8s-controller-demo/`):

```bash
docker build -f 06-finalizers/Dockerfile -t <your-registry>/trackedresource-controller:latest .
docker push <your-registry>/trackedresource-controller:latest
```

If you are currently in `06-finalizers/`, use parent directory (`..`) as build context:

```bash
docker build -f Dockerfile -t <your-registry>/trackedresource-controller:latest ..
docker push <your-registry>/trackedresource-controller:latest
```

### 2. Set your image in the Deployment manifest

Edit `config/manager/deployment.yaml` and replace:

```text
docker.io/your-user/trackedresource-controller:latest
```

with your pushed image.

### 3. Install CRD + RBAC + controller Deployment

```bash
kubectl apply -f config/crd/
kubectl apply -f config/rbac/
kubectl apply -f config/manager/
```

### 4. Create a sample custom resource

```bash
kubectl apply -f config/samples/
```

### 5. Verify finalizer behavior

```bash
kubectl logs -n default deploy/trackedresource-controller -f
kubectl get trackedresource my-tracked -o yaml
```

Check that `metadata.finalizers` contains `learn.example.com/cleanup`.

### 6. Delete and observe cleanup flow

```bash
kubectl delete trackedresource my-tracked
kubectl get trackedresources.learn.example.com
```

In logs you should see cleanup execution followed by finalizer removal.

---

## Congratulations

You've completed the learning path. Here's what you now know:

1. **Controller loop**: Reconcile is event-driven, idempotent, and edge-triggered.
2. **Status updates**: Use `r.Status().Update()` to write observed state.
3. **Creating resources**: Use `CreateOrUpdate` for idempotent child resource management.
4. **Owner references**: Link children to parents for automatic garbage collection.
5. **Watching anything**: Controllers can watch any Kubernetes resource, not just CRDs.
6. **Finalizers**: Register cleanup hooks that run before deletion.

**Where to go next:**

- Read the [controller-runtime docs](https://pkg.go.dev/sigs.k8s.io/controller-runtime)
- Explore [kubebuilder](https://book.kubebuilder.io/) for a framework that scaffolds much of this
- Look at real-world operators like [cert-manager](https://github.com/cert-manager/cert-manager) or [external-secrets](https://github.com/external-secrets/external-secrets)
