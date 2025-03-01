"use client";

import MetricChart from "@/components/MetricChart";

export default function MetricsPage() {
  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-6">Prometheus Metrics Dashboard</h1>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-2 gap-6">
        <MetricChart 
          metricName="tasks_processed_total" 
          title="Total Tasks Processed" 
          yAxisLabel="Tasks" 
        />

        <MetricChart 
          metricName="task_processing_time_seconds" 
          title="Task Processing Time (Seconds)" 
          yAxisLabel="Seconds" 
        />

        <MetricChart 
          metricName="task_retries_total" 
          title="Total Task Retries" 
          yAxisLabel="Retries" 
        />

        <MetricChart 
          metricName="pending_tasks_total" 
          title="Pending Tasks in Queue" 
          yAxisLabel="Count" 
        />
      </div>
    </div>
  );
}
