# Distributed Task Scheduler - Services

## Overview
This directory contains the core **backend services** that power the **Distributed Task Scheduler**. These services are responsible for **task scheduling, execution, monitoring, and data persistence**.

## **System Architecture (Flow Representation)**
```
+------------+       +------------+        +-------------+
|   Clients  |       | API Gateway |        |  Scheduler  |
+------------+       +------------+        +-------------+
| Web UI     |<----->| Handles API |<----->| Manages     |
| CLI        |       | & WebSockets|       | Task Queue  |
+------------+       +------------+        +-------------+
                                        |
                                        v
                                  +------------+
                                  |   Worker   |
                                  +------------+
                                  | Executes   |
                                  | Tasks      |
                                  | Publishes  |
                                  | Kafka Events|
                                  +------------+
                                        |
                                        v
+------------+       +------------+        +-------------+
|  Storage   |       |  Monitoring |        |    UI      |
+------------+       +------------+        +-------------+
| PostgreSQL |<----->| Prometheus  |<----->| WebSockets  |
| Redis      |       | Grafana     |       | Real-time   |
| Elastic    |       +------------+        | Updates     |
+------------+                              +-------------+
```

## **Getting Started**

### **Build & Run Locally**
To start all services locally using **Docker Compose**, run:

```sh
make build
make run
```

To stop services:
```sh
make clean
```

To run an individual service:
```sh
docker-compose up <service-name>
```

Example:
```sh
docker-compose up scheduler
```

---

## **Deployment**
### **Kubernetes Deployment**
To deploy the backend services using Kubernetes:
```sh
make k8s-deploy
```

To install the **Helm chart**:
```sh
make helm-install
```

### **CI/CD Pipeline**
- On each release, **Docker images** are built and pushed to **Docker Hub**.
- Helm charts are packaged and pushed automatically.
- To manually trigger a release:
```sh
make release
```

---

## **Monitoring & Observability**
- **Prometheus & Grafana** are used for system metrics.
- **Elasticsearch & Kibana** provide real-time log search and analytics.

To start monitoring:
```sh
make monitor
```

This will forward ports to access **Grafana, Prometheus, and Elasticsearch**.

---

## **Storage & Caching**
- **PostgreSQL** → Stores task details and execution history.
- **Redis** → Caches frequently accessed task data.
- **Elasticsearch** → Stores logs for real-time search and analytics.

---
