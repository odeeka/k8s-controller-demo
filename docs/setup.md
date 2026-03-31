# Environment Setup

This guide walks you through setting up a local Kubernetes cluster and verifying your Go toolchain before running the examples.

---

## 1. Install Go

Requires Go ≥ 1.22.

```bash
# Check your version
go version
# Should print: go version go1.22.x or higher
```

Download from https://go.dev/dl/ if needed.

---

## 2. Install kubectl

```bash
# macOS (Homebrew)
brew install kubectl

# Linux (direct download)
curl -LO "https://dl.k8s.io/release/$(curl -sL https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
chmod +x kubectl && sudo mv kubectl /usr/local/bin/

# Verify
kubectl version --client
```

---

## 3. Install kind (Kubernetes in Docker)

`kind` runs a Kubernetes cluster inside Docker — lightweight and fast for local development.

```bash
# macOS (Homebrew)
brew install kind

# Linux (direct download)
go install sigs.k8s.io/kind@latest

# Verify
kind version
```

> Alternative: [minikube](https://minikube.sigs.k8s.io/docs/start/) also works. Adjust `kind` commands accordingly.

---

## 4. Create a Local Cluster

```bash
# Create a cluster named 'learn'
kind create cluster --name learn

# Verify it's running
kubectl cluster-info --context kind-learn
kubectl get nodes
```

You should see one node in `Ready` state.

---

## 5. Download Go Dependencies

From the repository root:

```bash
go mod tidy
```

This downloads all dependencies and generates `go.sum`. You only need to do this once (or after pulling new commits).

---

## 6. Verify Everything Works

Install the CRD for step 01 and apply a sample:

```bash
kubectl apply -f 01-minimal-controller/config/crd/
kubectl apply -f 01-minimal-controller/config/samples/

# Verify the CRD was installed
kubectl get crd helloworlds.learn.example.com

# Verify the sample CR was created
kubectl get helloworlds
```

Then run the controller:

```bash
go run ./01-minimal-controller/
```

You should see log output like:
```
INFO  Starting HelloWorld controller...
INFO  Reconciling HelloWorld  {"name": "World", "resource": "default/my-first-controller"}
```

Press `Ctrl+C` to stop.

---

## Cleanup

```bash
# Delete the test cluster when done
kind delete cluster --name learn
```

---

## Permissions Note

When running locally, the controller uses your `~/.kube/config` (or the `KUBECONFIG` env var). For a local kind or minikube cluster, this gives you `cluster-admin` permissions by default — which is fine for learning.

In a real cluster, you would create a `ServiceAccount` with specific RBAC permissions. Each step's `config/` directory contains a `role.yaml` example showing what permissions the controller needs.
