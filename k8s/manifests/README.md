# Distributed Task Scheduler - Kubernetes Manifests

## Overview
This directory contains the **Kubernetes manifests** required to deploy the **Distributed Task Scheduler**. It includes **Deployments, Services, ConfigMaps, Persistent Volume Claims (PVCs), and an Ingress configuration**.

### Deployment Components
The manifests define:
- **Backend (Go API)**
- **UI (Next.js)**
- **Database (PostgreSQL)**
- **Cache (Redis)**
- **Message Queue (Kafka & Zookeeper)**
- **Search & Logging (Elasticsearch & Kibana)**
- **Monitoring (Prometheus & Grafana)**
- **Ingress (Nginx)** for external access

---

## Architecture
```
+--------------------+         +------------------+
|   User Requests   | ----->  |   Nginx Ingress  |
+--------------------+         +------------------+
                                       |
  +------------------------------------+-----------------------------------+
  |                 Backend Services                                      |
  +---------+----------+----------+------------+-------------------------+
  | Backend | UI       | Kafka    | Storage    | Monitoring              |
  | (API)   | (React)  | (Events) | (DB + Redis + Elastic) | (Prometheus) |
  +---------+----------+----------+------------+-------------------------+
```

- **Users interact via UI (`ui`)**
- **Ingress routes API requests (`backend`)**
- **Task data stored in PostgreSQL & Redis**
- **Kafka manages task event streaming**
- **Prometheus collects system metrics**
- **Elasticsearch indexes logs for searchability**

---

## Setup & Deployment

### 1. Apply Manifests
```sh
make k8s-deploy
```

### 2. Verify Deployments
```sh
kubectl get pods -n task-scheduler
kubectl get services -n task-scheduler
kubectl get ingress -n task-scheduler
```

### 3. Access UI & API
- **Frontend UI:** `https://dst.jougan.live`
- **API Gateway:** `https://api-dst.jougan.live`

---

## Kubernetes Components

### 1. API Backend
- **Deployment**: `backend.yaml`
- **Service**: `backend-service.yaml`
- **ConfigMap**: `backend-config.yaml`
- **Environment Variables**:
  - `DB_HOST`: PostgreSQL Service
  - `KAFKA_BROKERS`: Kafka Service
  - `REDIS_HOST`: Redis Service
  - `ELASTICSEARCH_URL`: Elasticsearch Service
  - `PROMETHEUS_URL`: Prometheus Service

### 2. UI Service
- **Deployment**: `ui.yaml`
- **Service**: `ui-service.yaml`
- **Env Variables**:
  - `BACKEND_URL`: `https://api-dst.jougan.live`

### 3. Database (PostgreSQL)
- **Deployment**: `postgres.yaml`
- **Service**: `postgres-service.yaml`
- **Persistent Volume Claim (PVC)**: `postgres-pvc.yaml`

### 4. Redis (Cache)
- **Deployment**: `redis.yaml`
- **Service**: `redis-service.yaml`

### 5. Kafka & Zookeeper
- **Deployment**: `kafka.yaml`, `zookeeper.yaml`
- **Service**: `kafka-service.yaml`, `zookeeper-service.yaml`

### 6. Elasticsearch & Kibana
- **Deployment**: `elasticsearch.yaml`, `kibana.yaml`
- **Service**: `elasticsearch-service.yaml`, `kibana-service.yaml`
- **Persistent Volume Claim (PVC)**: `elasticsearch-pvc.yaml`

### 7. Monitoring (Prometheus & Grafana)
- **Deployment**: `prometheus.yaml`, `grafana.yaml`
- **Service**: `prometheus-service.yaml`, `grafana-service.yaml`
- **Persistent Volume Claim (PVC)**: `grafana-pvc.yaml`
- **Prometheus ConfigMap**: `prometheus-config.yaml`

### 8. Ingress
- **Ingress Config**: `ingress.yaml`
- Routes:
  - `dst.jougan.live` → UI
  - `api-dst.jougan.live` → Backend

---

## Monitoring & Debugging

### 1. View Logs
```sh
kubectl logs -n task-scheduler -l app=backend
```

### 2. Check System Metrics
```sh
make monitor
```

### 3. Access Prometheus & Grafana
```sh
kubectl port-forward svc/grafana -n task-scheduler 3000:3000 &
kubectl port-forward svc/prometheus -n task-scheduler 9090:9090 &
```
- **Grafana:** `http://localhost:3000`
- **Prometheus:** `http://localhost:9090`

---

