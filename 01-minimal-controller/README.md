# Step 01 — Minimal Controller

**New concepts:** The reconciliation loop · The Manager · Watching a custom resource

---

## What This Step Does

This is the simplest possible controller. It:

1. Registers a custom resource type called `HelloWorld`
2. Watches for `HelloWorld` objects in the cluster
3. Every time one is created, updated, or deleted — it logs a message

The controller does **nothing else**. No status updates, no child resources, no side effects. Its only purpose is to show you what the reconciliation loop looks like.

---

## Key Ideas

### The Reconcile function

```go
func (r *HelloWorldReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error)
```

This function is called automatically by controller-runtime whenever a `HelloWorld` resource changes. You don't call it yourself — the framework does.

`req` contains only the **name and namespace** of the changed resource. The first thing you do is fetch the full resource with `r.Get(...)`.

### Event-driven, not polling

The framework watches the Kubernetes API for events (create, update, delete). When an event fires, it puts the resource's name into a queue, and a worker calls `Reconcile`. You don't need to poll.

### Returning from Reconcile

- `ctrl.Result{}, nil` → success, wait for the next event
- `ctrl.Result{}, err` → failure, requeue with exponential backoff
- `ctrl.Result{RequeueAfter: d}, nil` → requeue after a specific duration

---

## Project Layout

```text
01-minimal-controller/
├── main.go                          # Manager setup
├── api/v1/
│   ├── groupversion_info.go         # API group + scheme registration
│   ├── helloworld_types.go          # HelloWorld struct definition
│   └── zz_generated.deepcopy.go    # DeepCopy methods (normally code-generated)
├── internal/controller/
│   └── helloworld_controller.go     # The Reconciler
└── config/
    ├── crd/
    │   └── learn.example.com_helloworlds.yaml  # CRD manifest
    └── samples/
        └── helloworld_sample.yaml   # A sample HelloWorld object
```

---

## How to Run

```bash
# 1. Install the CRD into your cluster
kubectl apply -f config/crd/

# 2. Apply a sample HelloWorld resource
kubectl apply -f config/samples/

# 3. Run the controller (uses your current kubeconfig)
go run .
```

You should see output like:

```text
INFO  Reconciling HelloWorld  {"name": "World", "resource": "default/my-first-controller"}
```

## Run In-Cluster (Deployment)

If you want the controller to run *inside* Kubernetes instead of on your laptop, this step includes minimal RBAC + Deployment manifests under `config/rbac/` and `config/manager/`.

### 1. Build and push image

From the repository root (`k8s-controller-demo/`), build the image using the step-local Dockerfile:

```bash
docker build -f 01-minimal-controller/Dockerfile -t <your-registry>/helloworld-controller:latest .
docker push <your-registry>/helloworld-controller:latest
```

If you are currently in `01-minimal-controller/`, use parent directory (`..`) as build context:

```bash
docker build -f Dockerfile -t <your-registry>/helloworld-controller:latest ..
docker push <your-registry>/helloworld-controller:latest
```

### 2. Set your image in the Deployment manifest

Edit:

```text
config/manager/deployment.yaml
```

and replace:

```text
ghcr.io/your-user/helloworld-controller:latest
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

### 5. Check controller logs

```bash
kubectl logs -n default deploy/helloworld-controller -f
```

You should see reconcile logs for `my-first-controller`.

### 6. Trigger another reconcile

```bash
kubectl patch helloworld my-first-controller -n default --type=merge -p '{"spec":{"name":"ClusterRun"}}'
```

The Deployment pod should log a new reconciliation.

---

## Try It

While the controller is running, open another terminal and experiment:

```bash
# Create another HelloWorld — controller will reconcile it
kubectl apply -f - <<EOF
apiVersion: learn.example.com/v1
kind: HelloWorld
metadata:
  name: second
  namespace: default
spec:
  name: Kubernetes
EOF

# Update it — controller reconciles again
kubectl patch helloworld second --type=merge -p '{"spec":{"name":"Controller"}}'

# Delete it — controller reconciles one last time (and IgnoreNotFound handles the missing resource)
kubectl delete helloworld second
```

Notice: deleting the resource still triggers `Reconcile`, but `r.Get` returns "not found", and `client.IgnoreNotFound(err)` converts that to `nil`. The reconciler exits cleanly.

---

## What's Next

[Step 02 →](../02-status-updates/) Writing back to `.status` on the custom resource.
