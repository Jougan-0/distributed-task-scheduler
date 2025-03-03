"use client";

import { useEffect, useState } from "react";
import axios from "axios";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Legend,
} from "recharts";

const BACKEND_URL =
  process.env.NEXT_PUBLIC_BACKEND_URL || "http://localhost:8080";

interface MetricChartProps {
  promQuery: string;
  title: string;
  yAxisLabel: string;
  rangeSeconds?: number;
  stepSeconds?: number;
}

function parsePrometheusData(results: any[]): any[] {
  const dataMap: Record<string, Record<string, number | string>> = {};

  results.forEach((series) => {
    const labelString = Object.entries(series.metric)
      .map(([k, v]) => `${k}=${v}`)
      .join(", ");

    series.values.forEach(([timestamp, value]: [number, string]) => {
      const timeStr = new Date(timestamp * 1000).toLocaleTimeString();
      if (!dataMap[timeStr]) {
        dataMap[timeStr] = { time: timeStr };
      }
      dataMap[timeStr][labelString] = parseFloat(value);
    });
  });

  const allTimes = Object.keys(dataMap).sort((a, b) => {
    return (
      new Date(`1970-01-01T${a}`).getTime() -
      new Date(`1970-01-01T${b}`).getTime()
    );
  });

  return allTimes.map((t) => dataMap[t]);
}

export default function MetricChart({
  promQuery,
  title,
  yAxisLabel,
  rangeSeconds = 1800,
  stepSeconds = 15,
}: MetricChartProps) {
  const [data, setData] = useState<any[]>([]);
  const [seriesKeys, setSeriesKeys] = useState<string[]>([]);

  useEffect(() => {
    const fetchMetrics = async () => {
      try {
        const end = Math.floor(Date.now() / 1000);
        const start = end - rangeSeconds;
        const step = stepSeconds;

        const res = await axios.get(`${BACKEND_URL}/api/v1/prometheus/query`, {
          params: {
            query: promQuery,
            start,
            end,
            step,
          },
        });

        const result = res.data.data.result || [];
        if (result.length > 0) {
          const chartData = parsePrometheusData(result);
          setData(chartData);

          const allKeys = new Set<string>();
          chartData.forEach((item) => {
            Object.keys(item).forEach((k) => {
              if (k !== "time") {
                allKeys.add(k);
              }
            });
          });
          setSeriesKeys(Array.from(allKeys));
        } else {
          setData([]);
          setSeriesKeys([]);
        }
      } catch (err) {
        console.error(`Error fetching data for query "${promQuery}":`, err);
      }
    };

    fetchMetrics();
    const interval = setInterval(fetchMetrics, 5000);
    return () => clearInterval(interval);
  }, [promQuery, rangeSeconds, stepSeconds]);

  return (
    <div className="bg-white p-6 rounded-lg shadow">
      <h2 className="text-lg font-semibold mb-4">{title}</h2>
      <ResponsiveContainer width="100%" height={300}>
        <LineChart data={data}>
          <CartesianGrid strokeDasharray="3 3" stroke="#ccc" />
          <XAxis dataKey="time" stroke="#000" />
          <YAxis
            stroke="#000"
            label={{
              value: yAxisLabel,
              angle: -90,
              position: "insideLeft",
              fill: "#000",
            }}
          />
          <Tooltip
            contentStyle={{ backgroundColor: "#fff", color: "#000" }}
            labelStyle={{ color: "#000" }}
            itemStyle={{ color: "#000" }}
          />
          <Legend wrapperStyle={{ color: "#000" }} />

          {seriesKeys.map((key) => (
            <Line
              key={key}
              type="monotone"
              dataKey={key}
              strokeWidth={2}
              dot={false}
            />
          ))}
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}
