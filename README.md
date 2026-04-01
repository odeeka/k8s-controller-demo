# Kubernetes Controller Learning Path

A hands-on guide to writing Kubernetes controllers (operators) in Go using the **controller-runtime** library.

Each step is a self-contained, runnable example that introduces exactly one or two new ideas. Work through them in order to build a solid intuition for how controllers work — from first principles up to practical patterns.

---

## What You Will Learn

| Step | Description | New Concepts |
|------|-------------|--------------|
| [01 — Minimal Controller](./01-minimal-controller/) | Watch a CRD and log reconcile events | Controller loop, Manager |
| [02 — Status Updates](./02-status-updates/) | Write back to `.status` on a CR | Status subresource, idempotency |
| [03 — ConfigMap from CR](./03-configmap-from-cr/) | Create a ConfigMap from CR data | Creating child resources, idempotent creates |
| [04 — Deployment Manager](./04-deployment-manager/) | Manage a Deployment from a CR | Owner references, watching owned resources |
| [05 — Label Enforcer](./05-label-enforcer/) | Enforce labels on Namespaces | Watching built-in resources, patching |
| [06 — Finalizers](./06-finalizers/) | Clean up external state before deletion | Finalizers, deletion lifecycle |

---

## Prerequisites

- [Go](https://go.dev/dl/) ≥ 1.22
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- A local cluster: [kind](https://kind.sigs.k8s.io/) (recommended) or [minikube](https://minikube.sigs.k8s.io/)
- Basic familiarity with Kubernetes objects (Pod, Deployment, ConfigMap)

> See [docs/setup.md](./docs/setup.md) for cluster setup instructions.

---

## Quick Start

Install `go`

```bash
wget https://go.dev/dl/go1.26.1.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.26.1.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

Install `kind`

```bash
[$(uname -m) = x86_64 ] && curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.31.0/kind-$(uname)-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind
```

Bootstrap for repo and cluster:

```bash
# 1. Clone the repo
git clone https://github.com/example/k8s-controller-demo
cd k8s-controller-demo

# 2. Download Go dependencies
go mod tidy

# 3. Create a local cluster
kind create cluster --config ./kind/kind-cluster.yaml --name demo

# 4. Start with step 01
cd 01-minimal-controller
kubectl apply -f config/crd/
kubectl apply -f config/samples/
go run .
```

Destroy the cluster after learning:

```bash
kind delete cluster
```

---

## Project Structure

```text
k8s-controller-demo/
├── go.mod                        # Single Go module for all steps
├── docs/
│   ├── 00-concepts.md            # Background: controllers, reconciliation loops
│   └── setup.md                  # Cluster setup guide
├── 01-minimal-controller/        # Step 1: watch + log
├── 02-status-updates/            # Step 2: write .status
├── 03-configmap-from-cr/         # Step 3: create child resources
├── 04-deployment-manager/        # Step 4: owner references
├── 05-label-enforcer/            # Step 5: watch built-in resources
└── 06-finalizers/                # Step 6: cleanup on deletion
```

Each step directory contains:

- `README.md` — what this step teaches and how to run it
- `main.go` — manager setup and entrypoint
- `api/v1/` — custom resource type definitions
- `internal/controller/` — the reconciler logic
- `config/crd/` — CRD YAML to install into the cluster
- `config/samples/` — sample CR to test with

---

## Background Reading

Before diving in, read [docs/00-concepts.md](./docs/00-concepts.md) for a primer on:

- What a controller is
- How the reconciliation loop works
- The role of the Manager, Scheme, and Client

---

## Running a Step

Each step can be run from its own directory:

```bash
cd 01-minimal-controller

# Install the CRD
kubectl apply -f config/crd/

# Apply a sample resource
kubectl apply -f config/samples/

# Run the controller (uses your current kubeconfig)
go run .
```

The controller runs in the foreground. Press `Ctrl+C` to stop it.

> **Note:** The controller binds to `:8080` for metrics by default. Make sure this port is free, or see each step's README for how to change it.

---

## Learning Tips

- Read the code top-to-bottom — every file is commented to explain *why*, not just *what*.
- Try breaking things: delete a resource while the controller is running, change a spec field, etc.
- Each step builds on the previous one — resist the urge to skip ahead.
- The code is deliberately simple. Production controllers add more layers (error handling, conditions, leader election) — those patterns are introduced gradually.
