# Step 05 — Label Enforcer

**New concepts:** Watching built-in resources · Patch vs Update · RBAC for controllers

---

## What This Step Does

This step has **no CRD**. The controller watches `Namespace` objects and ensures every user-created namespace has a `team` label. If a namespace is missing the label, the controller adds `team: unassigned` automatically.

This demonstrates two important things:

1. Controllers don't need custom resources — they can watch *any* Kubernetes resource.
2. Using `r.Patch()` instead of `r.Update()` for partial changes.

---

## Key Ideas

### Watching built-in resources

You can watch any resource that is registered in the scheme. Since we call `clientgoscheme.AddToScheme(scheme)` in `main.go`, all standard Kubernetes types (including `Namespace`) are available.

```go
ctrl.NewControllerManagedBy(mgr).
    For(&corev1.Namespace{}).
    Complete(r)
```

### Patch vs Update

`r.Update()` replaces the full resource spec. If another actor has since modified the resource, you risk overwriting their changes with a stale version.

`r.Patch()` sends only the fields you want to change. This is safer for resources like Namespaces that may be managed by multiple parties.

```go
// MergeFrom creates a JSON Merge Patch from the base snapshot.
// Only fields that are different from the base are included in the patch.
patch := client.MergeFrom(ns.DeepCopy())   // snapshot the current state
ns.Labels["team"] = "unassigned"           // make the change
r.Patch(ctx, &ns, patch)                   // sends only {"metadata":{"labels":{"team":"unassigned"}}}
```

### Skipping system namespaces

The controller skips `kube-system`, `kube-public`, and `kube-node-lease` — these are managed by Kubernetes itself and should not be modified.

### RBAC

Because this controller modifies Namespaces (cluster-scoped resources), it needs a `ClusterRole` instead of a regular `Role`. See `config/rbac/clusterrole.yaml`.

---

## Project Layout

```text
05-label-enforcer/
├── main.go
├── internal/controller/
│   └── namespace_controller.go  # watches Namespaces, patches labels
└── config/
    ├── rbac/
    │   └── clusterrole.yaml     # ClusterRole for reading/updating Namespaces
    └── samples/
        └── test_namespace.yaml  # A sample namespace to test with
```

No `api/v1/` directory — there is no CRD.

---

## How to Run

```bash
# No CRD to install — just run straight away
go run .
```

In another terminal, create a namespace without a `team` label:

```bash
kubectl apply -f config/samples/test_namespace.yaml

# Check: the controller should add the label within seconds
kubectl get namespace learning-test --show-labels
# NAME            LABELS
# learning-test   kubernetes.io/metadata.name=learning-test,team=unassigned
```

Try adding a namespace with a label already set:

```bash
kubectl create namespace my-team-ns
kubectl label namespace my-team-ns team=backend
# Controller sees the label is already set and does nothing
```

## Run In-Cluster (Deployment)

If you want this controller to run inside Kubernetes, use the manifests under `config/rbac/` and `config/manager/`.

### 1. Build and push image

From repository root (`k8s-controller-demo/`):

```bash
docker build -f 05-label-enforcer/Dockerfile -t <your-registry>/namespace-label-controller:latest .
docker push <your-registry>/namespace-label-controller:latest
```

If you are currently in `05-label-enforcer/`, use parent directory (`..`) as build context:

```bash
docker build -f Dockerfile -t <your-registry>/namespace-label-controller:latest ..
docker push <your-registry>/namespace-label-controller:latest
```

### 2. Set your image in the Deployment manifest

Edit `config/manager/deployment.yaml` and replace:

```text
docker.io/your-user/namespace-label-controller:latest
```

with your pushed image.

### 3. Install RBAC + controller Deployment

```bash
kubectl apply -f config/rbac/
kubectl apply -f config/manager/
```

### 4. Test label enforcement

```bash
kubectl apply -f config/samples/
kubectl logs -n default deploy/namespace-label-controller -f
kubectl get namespace learning-test --show-labels
```

Check that the `team=unassigned` label appears automatically.

### 5. Test idempotent behavior

```bash
kubectl create namespace my-team-ns
kubectl label namespace my-team-ns team=backend
kubectl get namespace my-team-ns --show-labels
```

The controller should not overwrite an already set `team` label.

---

## What's Next

[Step 06 →](../06-finalizers/) Registering cleanup hooks that run before a resource is deleted.
