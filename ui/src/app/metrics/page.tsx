"use client";
export const dynamic = "force-dynamic";
export const revalidate = 0;
import MetricChart from "@/components/MetricChart";
export default function MetricsPage() {
  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-6 text-white">Prometheus Metrics Dashboard</h1>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-2 gap-6">
        <MetricChart
          promQuery="tasks_processed_total"
          title="Total Tasks Processed"
          yAxisLabel="Tasks"
        />

        <MetricChart
          promQuery="rate(task_processing_time_seconds_sum[5m]) / rate(task_processing_time_seconds_count[5m])"
          title="Avg Task Processing Time (last 5m)"
          yAxisLabel="Seconds"
        />

        <MetricChart
          promQuery="sum by (task_type) (task_retries_total)"
          title="Task Retries by Type"
          yAxisLabel="Retries"
        />

        <MetricChart
          promQuery="pending_tasks_total"
          title="Pending Tasks in Queue"
          yAxisLabel="Count"
        />
      </div>
    </div>
  );
}
