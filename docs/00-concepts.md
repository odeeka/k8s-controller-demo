# Kubernetes Controllers — Core Concepts

This document explains the foundational ideas behind Kubernetes controllers before you write any code. Read it once, then come back when something in the code examples isn't clear.

---

## What Is a Controller?

A Kubernetes **controller** is a program that watches the state of your cluster and takes actions to move it toward a desired state.

Every built-in Kubernetes feature is implemented as a controller:
- The **Deployment controller** watches Deployment objects and creates/updates ReplicaSets.
- The **ReplicaSet controller** watches ReplicaSets and creates/deletes Pods.
- The **Job controller** watches Jobs and creates Pods to run work.

When you write a custom controller (also called an **operator**), you follow the same pattern — just with your own resource types and logic.

---

## The Reconciliation Loop

The core of every controller is the **reconciliation loop**:

```
┌─────────────────────────────────────────────────┐
│                                                 │
│  ┌──────────┐    watch     ┌──────────────────┐ │
│  │          │◄─────────────│  Kubernetes API  │ │
│  │Reconciler│              │  (etcd)          │ │
│  │          │──────────────►                  │ │
│  └──────────┘    act       └──────────────────┘ │
│                                                 │
└─────────────────────────────────────────────────┘
```

1. **Watch**: The controller watches for changes to resources (creates, updates, deletes).
2. **Reconcile**: When a change is detected, the `Reconcile` function is called with the name/namespace of the changed resource.
3. **Act**: The function reads the current state, compares it to the desired state, and makes API calls to close the gap.
4. **Repeat**: The loop runs continuously, re-triggered by new events.

### The `Reconcile` function — key properties

```go
func (r *MyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error)
```

- `req` contains only the **name and namespace** of the resource that changed, not the resource itself. You must fetch it with `r.Get(...)`.
- The function should be **idempotent**: running it 10 times on the same state should have the same result as running it once.
- Return `ctrl.Result{}, nil` to say "I'm done, wait for the next event."
- Return `ctrl.Result{}, err` to requeue immediately (with backoff).
- Return `ctrl.Result{RequeueAfter: d}, nil` to requeue after a delay.

---

## Desired State vs. Observed State

Kubernetes resources have two parts:

| Part | Field | Who writes it? | Meaning |
|------|-------|----------------|---------|
| **Spec** | `.spec` | You (the user) | What you *want* |
| **Status** | `.status` | The controller | What *is* (observed) |

Your controller's job is to look at `.spec`, check what *actually* exists in the cluster, and take actions to make reality match the spec. Then it writes back to `.status` to record what it observed.

---

## The Manager

The **Manager** is controller-runtime's main object. It handles:

- Connecting to the Kubernetes API server
- Running multiple controllers in a single process
- Caching API responses (to avoid hammering the API)
- Leader election (so only one instance is active at a time — not used in these examples)
- Propagating signals (Ctrl+C → graceful shutdown)

```go
mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
    Scheme: scheme,
})
```

---

## The Scheme

The **Scheme** is a registry that maps Go types to Kubernetes API group/version/kind strings.

```go
// This tells the runtime: "when you see a JSON object with
// apiVersion: apps/v1 and kind: Deployment, decode it into appsv1.Deployment"
utilruntime.Must(clientgoscheme.AddToScheme(scheme))

// And: "when you see apiVersion: learn.example.com/v1 and kind: HelloWorld,
// decode it into learnv1.HelloWorld"
utilruntime.Must(learnv1.AddToScheme(scheme))
```

You must register every type you want to work with, both standard Kubernetes types and your own custom types.

---

## The Client

Inside a controller, `r.Client` (or just `r.Get`, `r.Create`, etc.) is how you talk to the Kubernetes API:

```go
// Read
r.Get(ctx, req.NamespacedName, &myResource)

// Create
r.Create(ctx, &newConfigMap)

// Update
r.Update(ctx, &existingResource)

// Update only the status (requires status subresource)
r.Status().Update(ctx, &resource)

// Delete
r.Delete(ctx, &resource)

// Patch (partial update, safer than full Update)
r.Patch(ctx, &resource, patch)
```

The client reads from the **cache** by default (fast, local), and writes go directly to the **API server**.

---

## Custom Resource Definitions (CRDs)

A **CRD** registers a new resource type with Kubernetes. Once installed, you can create instances of it like any other resource.

```yaml
apiVersion: learn.example.com/v1
kind: HelloWorld
metadata:
  name: my-resource
spec:
  name: World
```

The CRD YAML file tells Kubernetes the schema (what fields are allowed, which are required). Your controller then watches for objects of this type.

---

## Key Terms Quick Reference

| Term | Meaning |
|------|---------|
| **Controller** | A loop that watches resources and takes action |
| **Operator** | A controller that manages a complex application |
| **Reconciler** | The struct that implements the `Reconcile` method |
| **CRD** | CustomResourceDefinition — registers a new resource type |
| **CR** | Custom Resource — an instance of a CRD |
| **Manager** | Runs controllers, handles caching and signals |
| **Scheme** | Maps Go types to API group/version/kind |
| **Finalizer** | Prevents deletion until cleanup is done |
| **Owner reference** | Links a child resource to its parent |
