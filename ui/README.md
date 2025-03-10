# **Distributed Task Scheduler - UI**

## **Overview**
This is the **frontend interface** for the **Distributed Task Scheduler**. It provides **real-time task monitoring, Kafka event tracking, Redis key inspection, and Prometheus metrics visualization**.

### **Key Features**
- **Task Management** → Create, View, and Search tasks.
- **Kafka Event Streaming** → Monitors task lifecycle (`TaskCreated`, `TaskCompleted`, `TaskFailed`).
- **Real-time Logs** → WebSocket-based live log streaming.
- **Prometheus Metrics** → Fetches and visualizes Prometheus metrics.
- **Redis Key Inspection** → Monitors task cache values.
- **Elasticsearch Search** → Queries stored task data.

---

## **Architecture**
```
+-----------------+        +--------------------+
|  User (Web UI)  | -----> |  API Gateway       |
+-----------------+        |  (REST + WebSocket)|
                           +--------------------+
                                    |
  +---------------------------------+---------------------------------+
  |                 Backend Services                                |
  +---------+----------+----------+------------+--------------------+
  | Scheduler | Worker | Kafka    | Storage    | Monitoring        |
  | (Tasks)   | (Exec) | (Events) | (DB + Redis + Elastic) | (Metrics) |
  +---------+----------+----------+------------+--------------------+
```
- The **UI communicates with the API Gateway** to interact with **tasks, logs, and metrics**.
- **WebSockets** are used for **real-time logs**.
- **Kafka events** are fetched periodically to update task statuses.
- **Prometheus metrics** are visualized using **interactive charts**.

---

## **Setup & Running Locally**
### **1. Install Dependencies**
```sh
npm install
```
### **2. Start the UI**
```sh
make ui-run
```
This will start the frontend at **`http://localhost:3001`**.

---

## **WebSocket Integration**
- The UI **connects to the backend WebSocket** (`/ws`) to receive live log updates.
- Logs are displayed in the `LogsPage.tsx` component.

---

## **Prometheus Metrics**
- The **`MetricChart.tsx`** component fetches **Prometheus data** from `/api/v1/query_range`.
- Metrics include:
  - **Total Tasks Processed**
  - **Task Processing Time**
  - **Pending Task Count**
  - **Task Retries by Type**

---

## **Kafka Event Streaming**
- **`KafkaEventsPage.tsx`** fetches Kafka events from:
  ```sh
  GET /kafka/events
  ```
- Events include:
  - **Task Created**
  - **Task Completed**
  - **Task Failed**

---

## **Redis Data Inspection**
- **`RedisPage.tsx`** queries:
  ```sh
  GET /redis/keys
  ```
- Displays key-value pairs stored in **Redis**.

---

## **Searching Tasks in Elasticsearch**
- **`ElasticsearchPage.tsx`** queries Elasticsearch:
  ```sh
  GET /api/v1/tasks/search/{query}
  ```
- Allows full-text search of stored tasks.

---

## **Deployment**
- The UI is deployed on **AWS EC2** with a **custom domain**.
- **Helm is used** for deployment:
  ```sh
  make helm-ui-install
  ```

---