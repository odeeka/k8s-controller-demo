# ============================================================
# k8s-controller-demo — Makefile helpers
# ============================================================
# Usage: make <target>
# All `run-*` targets assume a valid KUBECONFIG is set.
# ============================================================

.PHONY: help tidy cluster delete-cluster \
        install-crds-01 install-crds-02 install-crds-03 \
        install-crds-04 install-crds-06 \
        run-01 run-02 run-03 run-04 run-05 run-06

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ------------------------------------------------------------
# Dependency management
# ------------------------------------------------------------
tidy: ## Download and tidy Go dependencies (generates go.sum)
	go mod tidy

# ------------------------------------------------------------
# Local cluster (kind)
# ------------------------------------------------------------
cluster: ## Create a local kind cluster named 'learn'
	kind create cluster --name learn

delete-cluster: ## Delete the local kind cluster
	kind delete cluster --name learn

# ------------------------------------------------------------
# Install CRDs into the current cluster
# ------------------------------------------------------------
install-crds-01: ## Install CRDs for step 01
	kubectl apply -f 01-minimal-controller/config/crd/

install-crds-02: ## Install CRDs for step 02
	kubectl apply -f 02-status-updates/config/crd/

install-crds-03: ## Install CRDs for step 03
	kubectl apply -f 03-configmap-from-cr/config/crd/

install-crds-04: ## Install CRDs for step 04
	kubectl apply -f 04-deployment-manager/config/crd/

install-crds-06: ## Install CRDs for step 06
	kubectl apply -f 06-finalizers/config/crd/

# ------------------------------------------------------------
# Run each step
# ------------------------------------------------------------
run-01: ## Run step 01 — Minimal Controller
	go run ./01-minimal-controller/

run-02: ## Run step 02 — Status Updates
	go run ./02-status-updates/

run-03: ## Run step 03 — ConfigMap from CR
	go run ./03-configmap-from-cr/

run-04: ## Run step 04 — Deployment Manager
	go run ./04-deployment-manager/

run-05: ## Run step 05 — Label Enforcer (no CRD)
	go run ./05-label-enforcer/

run-06: ## Run step 06 — Finalizers
	go run ./06-finalizers/
