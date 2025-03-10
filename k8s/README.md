# Distributed Task Scheduler - Kubernetes Overview

## Overview
This directory (`k8s/`) contains all the necessary **Kubernetes configurations** for deploying the **Distributed Task Scheduler**. It is structured into two subdirectories:

1. **`k8s/manifest/`** → Contains raw Kubernetes YAML manifests.
2. **`k8s/helm/`** → Contains a Helm chart for streamlined deployment.

Both approaches allow deploying **Backend, UI, Database, Redis, Kafka, Elasticsearch, Prometheus, and Grafana**.

---

## Directory Structure
```
k8s/
│── manifest/       # Kubernetes raw YAML manifests
│── helm/           # Helm chart for deploying to Kubernetes
```

---

## Deployment Options

### 1. Deploy using **Raw Manifests**
```sh
make k8s-deploy
```
- This applies all the **Kubernetes YAML files** from `k8s/manifest/`.

### 2. Deploy using **Helm**
```sh
make helm-install
```
- This installs the **Helm chart** from `k8s/helm/`.

---

## Key Components

### ✅ **API Backend**
- **Deployment:** `backend.yaml`
- **Service:** `backend-service.yaml`
- **ConfigMap:** `backend-config.yaml`

### ✅ **UI Service**
- **Deployment:** `ui.yaml`
- **Service:** `ui-service.yaml`

### ✅ **Database & Cache**
- **PostgreSQL Deployment & Service**
- **Redis Deployment & Service**

### ✅ **Message Queue & Event Streaming**
- **Kafka & Zookeeper Deployments**

### ✅ **Monitoring & Logging**
- **Prometheus & Grafana for Metrics**
- **Elasticsearch & Kibana for Log Search**

### ✅ **Ingress & Networking**
- **Nginx Ingress Controller**
- **Ingress Routes for API & UI**

---

## CI/CD Pipeline
- On **each release**, **Docker images are built and pushed** to **Docker Hub**.
- The **Helm chart is updated automatically** to deploy the latest version.

---

## Monitoring & Debugging

### 1. Check running services
```sh
kubectl get pods -n task-scheduler
kubectl get services -n task-scheduler
```

### 2. View logs
```sh
kubectl logs -n task-scheduler -l app=backend
```

### 3. Access Prometheus & Grafana
```sh
kubectl port-forward svc/grafana -n task-scheduler 3000:3000 &
kubectl port-forward svc/prometheus -n task-scheduler 9090:9090 &
```
- **Grafana:** `http://localhost:3000`
- **Prometheus:** `http://localhost:9090`

---
