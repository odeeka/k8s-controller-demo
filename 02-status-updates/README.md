# Step 02 — Status Updates

**New concepts:** The status subresource · Writing observed state · Idempotent updates

---

## What This Step Does

We introduce a `Greeter` custom resource with a proper `.status` block. The controller:

1. Reads `spec.greeting` and `spec.targetName`
2. Computes a greeting message
3. Writes `status.phase` and `status.message` back to the resource

This is the fundamental loop: **read spec → compute desired state → update status**.

---

## Key Ideas

### Status is for the controller, spec is for the user

| Field | Written by | Meaning |
|-------|-----------|---------|
| `.spec` | You (kubectl, GitOps) | What you *want* |
| `.status` | The controller | What the controller *observed/did* |

Never put controller-computed data in `.spec`. Never let users modify `.status` directly.

### The status subresource

Adding `// +kubebuilder:subresource:status` to your type (and `subresources: status: {}` in the CRD) enables the **status subresource**. This means:

- A normal `r.Update()` call ignores `.status` changes — they won't be persisted.
- You **must** use `r.Status().Update()` to update status fields.

This separation prevents accidental overwriting of status when a user edits the spec.

### Idempotency

The controller checks whether the status already matches what it would write. If nothing changed, it returns early without making an API call. This is called **idempotency**: applying the same reconcile repeatedly produces the same result.

```go
if gr.Status.Phase == "Ready" && gr.Status.Message == desiredMessage {
    log.Info("already up to date, skipping")
    return ctrl.Result{}, nil
}
```

This is a good habit. Without it, you'd generate spurious update events that re-trigger reconciliation unnecessarily.

---

## Project Layout

```
02-status-updates/
├── main.go
├── api/v1/
│   ├── groupversion_info.go
│   ├── greeter_types.go          # ← Now has a Status with Phase + Message
│   └── zz_generated.deepcopy.go
├── internal/controller/
│   └── greeter_controller.go     # ← Uses r.Status().Update()
└── config/
    ├── crd/learn.example.com_greeters.yaml
    └── samples/greeter_sample.yaml
```

---

## How to Run

```bash
kubectl apply -f config/crd/
kubectl apply -f config/samples/
go run .
```

After the controller starts, check the status:

```bash
kubectl get greeter my-greeter -o yaml
# Look for .status.phase and .status.message
```

Or use the printer columns:
```bash
kubectl get greeters
# NAME        PHASE   MESSAGE
# my-greeter  Ready   Hello, Kubernetes!
```

---

## Try It

```bash
# Change the greeting — controller reconciles and updates status
kubectl patch greeter my-greeter --type=merge -p '{"spec":{"greeting":"Howdy","targetName":"Partner"}}'

# Check the updated status
kubectl get greeter my-greeter -o jsonpath='{.status.message}'
```

---

## What's Next

[Step 03 →](../03-configmap-from-cr/) Creating a child resource (ConfigMap) from the spec.
