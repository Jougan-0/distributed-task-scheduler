"use client";

import { useEffect, useState } from "react";
import axios from "axios";
import {
  LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip,
  ResponsiveContainer, Legend,
} from "recharts";

const PROMETHEUS_URL = process.env.NEXT_PUBLIC_PROMETHEUS_URL || "http://localhost:9090";

interface MetricChartProps {
  metricName: string;
  title: string;
  yAxisLabel: string;
}

export default function MetricChart({ metricName, title, yAxisLabel }: MetricChartProps) {
  const [data, setData] = useState<any[]>([]);

  useEffect(() => {
    const fetchMetrics = async () => {
      try {
        const end = Math.floor(Date.now() / 1000);
        const start = end - 1800; 
        const step = 15;

        const res = await axios.get(`${PROMETHEUS_URL}/api/v1/query_range`, {
          params: {
            query: metricName,
            start,
            end,
            step,
          },
        });

        if (res.data.data.result.length > 0) {
          const metricData = res.data.data.result[0].values.map(
            ([timestamp, value]: [number, string]) => ({
              time: new Date(timestamp * 1000).toLocaleTimeString(),
              value: parseFloat(value),
            })
          );
          setData(metricData);
        } else {
          setData([]);
        }
      } catch (err) {
        console.error(`Error fetching ${metricName} data:`, err);
      }
    };

    fetchMetrics();
    const interval = setInterval(fetchMetrics, 5000);
    return () => clearInterval(interval);
  }, [metricName]);

  return (
    <div className="bg-white p-6 rounded-lg shadow">
      <h2 className="text-lg font-semibold mb-4">{title}</h2>
      <ResponsiveContainer width="100%" height={300}>
        <LineChart data={data}>
          <CartesianGrid strokeDasharray="3 3" stroke="#ccc" />
          <XAxis dataKey="time" stroke="#000" />
          <YAxis
            stroke="#000"
            label={{ value: yAxisLabel, angle: -90, position: "insideLeft", fill: "#000" }}
          />
          <Tooltip
            contentStyle={{ backgroundColor: "#fff", color: "#000" }}
            labelStyle={{ color: "#000" }}
            itemStyle={{ color: "#000" }}
          />
          <Legend wrapperStyle={{ color: "#000" }} />
          <Line
            type="monotone"
            dataKey="value"
            stroke="#FF0000"
            strokeWidth={2}
            dot={false}
          />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}
