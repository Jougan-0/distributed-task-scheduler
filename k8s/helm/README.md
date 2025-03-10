# Distributed Task Scheduler - Helm Deployment

## Overview
This directory contains the **Helm chart** for deploying the **Distributed Task Scheduler** on **Kubernetes**. The chart manages **Deployments, Services, ConfigMaps, Ingress, and Persistent Volume Claims (PVCs)**.

### Why Helm?
- **Automated Deployment** → Deploy and update the stack with a single command.
- **Configurability** → Customize environment variables and replicas via `values.yaml`.
- **CI/CD Integration** → Helm charts are automatically updated on each release.

---

## Helm Chart Structure
```
helm/
│── Chart.yaml         # Helm chart metadata
│── values.yaml        # Default configurations
│── values.prod.yaml   # Production configurations
│── templates/         # Kubernetes manifests (Helm templated)
│   ├── backend.yaml
│   ├── configmap.yaml
│   ├── elasticsearch.yaml
│   ├── grafana.yaml
│   ├── ingress.yaml
│   ├── kafka.yaml
│   ├── kibana.yaml
│   ├── postgres.yaml
│   ├── prometheus.yaml
│   ├── pvc.yaml
│   ├── redis.yaml
│   ├── secrets.yaml
│   ├── ui.yaml
│   ├── zookeeper.yaml
```

---

## Deployment

### 1. Install the Helm Chart
```sh
make helm-install
```
This installs the **Distributed Task Scheduler** with default values.

### 2. Upgrade Deployment
If you update the **Helm chart** (e.g., new image versions), run:
```sh
make helm-upgrade
```

### 3. Uninstall the Helm Chart
```sh
make helm-uninstall
```

---

## Configuring Helm Deployment
Modify **`values.yaml`** to customize settings like:
```yaml
backend:
  image: "shlok08/distributed-task-scheduler:v1.1"
  replicas: 2
  env:
    DB_HOST: "postgres.task-scheduler.svc.cluster.local"
    KAFKA_BROKERS: "kafka-broker.task-scheduler.svc.cluster.local:9092"
    REDIS_HOST: "redis.task-scheduler.svc.cluster.local"
    ELASTICSEARCH_URL: "http://elasticsearch.task-scheduler.svc.cluster.local:9200"
    PROMETHEUS_URL: "http://prometheus.task-scheduler.svc.cluster.local:9090"

ui:
  image: "shlok08/task-scheduler-ui:v2"
  replicas: 1
  env:
    BACKEND_URL: "http://api-dst.jougan.live"
```

To apply changes:
```sh
helm upgrade task-scheduler k8s/helm/ --values values.yaml
```

---

## CI/CD - Automatic Helm Chart Updates
- **GitHub Actions Pipeline** builds and pushes **Docker images** to **Docker Hub** on each release.
- The Helm chart is **packaged and updated automatically**.
- New releases deploy the latest **backend and UI images**.

---

## Monitoring & Debugging
### 1. Check Helm Releases
```sh
helm list -n task-scheduler
```

### 2. Inspect Running Pods
```sh
kubectl get pods -n task-scheduler
```

### 3. View Logs
```sh
kubectl logs -n task-scheduler -l app=backend
```

### 4. Check Service Endpoints
```sh
kubectl get services -n task-scheduler
```