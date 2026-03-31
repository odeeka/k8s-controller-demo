# Step 03 — ConfigMap from CR

**New concepts:** Creating child resources · Idempotent create-or-update · Reconciling against real cluster state

---

## What This Step Does

We introduce a `ConfigSource` custom resource with a `spec.data` field (a key-value map). The controller ensures a corresponding `ConfigMap` always exists in the cluster and stays in sync with the CR's data.

This demonstrates the core CRUD pattern for child resources:

```
Reconcile:
  1. Fetch the ConfigSource CR
  2. Check if the target ConfigMap already exists
     a. Not found  → Create it
     b. Found, data differs → Update it
     c. Found, data matches → Do nothing (idempotency)
  3. Update status to record what happened
```

---

## Key Ideas

### Creating child resources

Use `r.Create(ctx, &cm)` to create a new resource. Always check for the "already exists" error — another controller or the user may have created it first.

### Idempotent create-or-update

The pattern:
```go
var existing corev1.ConfigMap
err := r.Get(ctx, key, &existing)
if apierrors.IsNotFound(err) {
    return ctrl.Result{}, r.Create(ctx, desired)
} else if err != nil {
    return ctrl.Result{}, err
}
// existing found — update if needed
```

This is so common it has a helper in controller-runtime: `controllerutil.CreateOrUpdate`. We use it here to keep the code short and clear.

### Why no owner reference yet?

Owner references (making the ConfigMap automatically deleted when the ConfigSource is deleted) are introduced in **Step 04**. For now we deliberately omit them so you notice what happens:

```bash
kubectl delete configsource my-config
# ConfigMap still exists! This is a resource leak.
```

Step 04 fixes this with `controllerutil.SetControllerReference`.

---

## Project Layout

```
03-configmap-from-cr/
├── main.go
├── api/v1/
│   ├── groupversion_info.go
│   ├── configsource_types.go     # ← Has spec.data (map[string]string)
│   └── zz_generated.deepcopy.go # ← Map deep-copy example
├── internal/controller/
│   └── configsource_controller.go  # ← Uses CreateOrUpdate
└── config/
    ├── crd/learn.example.com_configsources.yaml
    └── samples/configsource_sample.yaml
```

---

## How to Run

```bash
kubectl apply -f config/crd/
kubectl apply -f config/samples/
go run .
```

Observe the created ConfigMap:
```bash
kubectl get configmap my-app-config -o yaml
```

Update the CR's data and watch the ConfigMap update:
```bash
kubectl patch configsource my-app-config --type=merge \
  -p '{"spec":{"data":{"key1":"updated-value"}}}'

kubectl get configmap my-app-config -o yaml
# key1 is now "updated-value"
```

---

## What's Next

[Step 04 →](../04-deployment-manager/) Owner references — so child resources are cleaned up automatically.
