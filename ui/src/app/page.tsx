export const dynamic = "force-dynamic";
"use client";
import Link from "next/link";
import { useState, useRef, useEffect } from "react";
export default function Home() {
  const [message, setMessage] = useState("");

  const [taskEmail, setTaskEmail] = useState("demo@example.com");
  const [taskType, setTaskType] = useState<"EMAIL" | "REPORT_GENERATION">("EMAIL");
  const [priority, setPriority] = useState<number>(1);

  const [showCustomForm, setShowCustomForm] = useState(false);

  const buttonRef = useRef<HTMLButtonElement>(null);
  const formRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    function handleClickOutside(e: MouseEvent) {
      if (showCustomForm) {
        if (
          formRef.current &&
          !formRef.current.contains(e.target as Node) &&
          buttonRef.current &&
          !buttonRef.current.contains(e.target as Node)
        ) {
          setShowCustomForm(false);
        }
      }
    }

    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [showCustomForm]);
  const BACKEND_URL = process.env.NEXT_PUBLIC_BACKEND_URL;
  const createExampleTask = async () => {
    setMessage("Creating example task...");
    try {
      const res = await fetch(`${BACKEND_URL}/api/v1/tasks`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          name: "DemoTask",
          type: "EMAIL",
          payload: JSON.stringify({ email: "demo@example.com" }),
          max_retries: 3,
          scheduled_time: "",
          priority: 1,
        }),
      });
      if (!res.ok) {
        throw new Error("Failed to create example task");
      }
      setMessage("DemoTask created successfully! Check logs/metrics pages.");
    } catch (err) {
      setMessage("Error creating example task: " + err);
    }
  };

  const createCustomTask = async () => {
    setMessage("Creating custom task...");
    try {
      const res = await fetch("http://localhost:8080/api/v1/tasks", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          name: `UserDefinedTask-${Date.now()}`,
          type: taskType,
          payload: JSON.stringify({ email: taskEmail }),
          max_retries: 3,
          scheduled_time: "",
          priority,
        }),
      });
      if (!res.ok) {
        throw new Error("Failed to create custom task");
      }
      setMessage(`Custom Task created successfully! (Type: ${taskType}, Email: ${taskEmail}, Priority: ${priority}). Check logs/metrics pages.`);
    } catch (err) {
      setMessage("Error creating custom task: " + err);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col">
      <div className="bg-gradient-to-r from-blue-600 to-purple-600 text-white py-16 px-8 flex flex-col items-center">
        <h1 className="text-4xl font-extrabold mb-4">Distributed Task Scheduler Dashboard</h1>
        <p className="text-lg max-w-2xl text-center">
          A unified platform to monitor and manage backend logs, Kafka events, Redis keys, Prometheus metrics, and Elasticsearch data—all in one place.
        </p>

        <div className="mt-6 flex gap-4">
          <Link
            href="/metrics"
            className="inline-block bg-white text-blue-600 font-semibold px-5 py-3 rounded-md shadow hover:bg-blue-100 transition"
          >
            View Metrics
          </Link>
          <button
            onClick={createExampleTask}
            className="inline-block bg-blue-100 text-blue-800 font-semibold px-5 py-3 rounded-md shadow hover:bg-blue-200 transition"
          >
            Create Example Task
          </button>
          <button
            onClick={() => setShowCustomForm(!showCustomForm)}
            ref={buttonRef}
            className="inline-block bg-blue-300 text-blue-900 font-semibold px-5 py-3 rounded-md shadow hover:bg-blue-400 transition"
          >
            {showCustomForm ? "Hide Custom Form" : "Create Custom Task"}
          </button>
        </div>

        {message && (
          <div className="mt-4 bg-white text-blue-800 py-2 px-4 rounded-md shadow">
            {message}
          </div>
        )}

        {showCustomForm && (
          <div
            ref={formRef}
            className="bg-white text-blue-800 py-4 px-6 rounded-md mt-8 shadow w-full max-w-xl"
          >
            <h2 className="text-xl font-bold mb-3">Create a Custom Task</h2>
            <div className="flex flex-col sm:flex-row gap-2 mb-4">
              <div className="flex-1">
                <label className="block font-semibold mb-1">Task Type</label>
                <select
                  value={taskType}
                  onChange={(e) =>
                    setTaskType(e.target.value as "EMAIL" | "REPORT_GENERATION")
                  }
                  className="border border-gray-300 rounded px-2 py-1 w-full"
                >
                  <option value="EMAIL">EMAIL</option>
                  <option value="REPORT_GENERATION">REPORT_GENERATION</option>
                </select>
              </div>

              <div className="flex-1">
                <label className="block font-semibold mb-1">Email (Payload)</label>
                <input
                  type="text"
                  value={taskEmail}
                  onChange={(e) => setTaskEmail(e.target.value)}
                  className="border border-gray-300 rounded px-2 py-1 w-full"
                />
              </div>

              <div className="flex-1">
                <label className="block font-semibold mb-1">Priority</label>
                <input
                  type="number"
                  value={priority}
                  onChange={(e) => setPriority(Number(e.target.value))}
                  className="border border-gray-300 rounded px-2 py-1 w-full"
                />
              </div>
            </div>
            <button
              onClick={createCustomTask}
              className="bg-blue-600 text-white px-4 py-2 rounded shadow hover:bg-blue-500 transition"
            >
              Create Custom Task
            </button>
          </div>
        )}
      </div>

      <div className="container mx-auto py-12 px-4 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 -mt-12">
        <Link href="/logs" className="bg-white rounded-xl p-6 shadow hover:shadow-lg transition">
          <h2 className="text-xl  font-bold mb-3 flex items-center gap-2">
            <svg
              className="w-6 h-6 text-gray-600"
              fill="none"
              stroke="currentColor"
              strokeWidth="1.5"
              viewBox="0 0 24 24"
            >
              <path strokeLinecap="round" strokeLinejoin="round" d="M3.75 7.5l16.5-4.125v17.25L3.75 18.75V7.5z" />
            </svg>
            Logs
          </h2>
          <p className="text-gray-700">
            View real-time backend logs streamed via WebSockets. Perfect for debugging and live monitoring.
          </p>
        </Link>

        <Link href="/kafka" className="bg-white rounded-xl p-6 shadow hover:shadow-lg transition">
          <h2 className="text-xl font-bold mb-3 flex items-center gap-2 ">
            <svg
              className="w-6 h-6 text-gray-600"
              fill="none"
              stroke="currentColor"
              strokeWidth="1.5"
              viewBox="0 0 24 24"
            >
              <path strokeLinecap="round" strokeLinejoin="round" d="M13.5 4.5L21 12m0 0l-7.5 7.5M21 12H3" />
            </svg>
            Kafka
          </h2>
          <p className="text-gray-700">
            Monitor Kafka events published by the task scheduler for deeper insights and async flows.
          </p>
        </Link>

        <Link href="/redis" className="bg-white rounded-xl p-6 shadow hover:shadow-lg transition">
          <h2 className="text-xl font-bold mb-3 flex items-center gap-2">
            <svg
              className="w-6 h-6 text-gray-600"
              fill="none"
              stroke="currentColor"
              strokeWidth="1.5"
              viewBox="0 0 24 24"
            >
              <path strokeLinecap="round" strokeLinejoin="round" d="M2.25 12.75l8.25-4.5 8.25 4.5-8.25 4.5-8.25-4.5z" />
            </svg>
            Redis
          </h2>
          <p className="text-gray-700">
            Inspect real-time Redis keys and data for fast lookups and caching verification.
          </p>
        </Link>

        <Link href="/metrics" className="bg-white rounded-xl p-6 shadow hover:shadow-lg transition">
          <h2 className="text-xl font-bold mb-3 flex items-center gap-2">
            <svg
              className="w-6 h-6 text-gray-600"
              fill="none"
              stroke="currentColor"
              strokeWidth="1.5"
              viewBox="0 0 24 24"
            >
              <path strokeLinecap="round" strokeLinejoin="round" d="M3 3v18h18" />
            </svg>
            Prometheus
          </h2>
          <p className="text-gray-700">
            Visualize task metrics like processed tasks, latency histograms, and more with interactive charts.
          </p>
        </Link>

        <Link href="/elasticsearch" className="bg-white rounded-xl p-6 shadow hover:shadow-lg transition">
          <h2 className="text-xl font-bold mb-3 flex items-center gap-2">
            <svg
              className="w-6 h-6 text-gray-600"
              fill="none"
              stroke="currentColor"
              strokeWidth="1.5"
              viewBox="0 0 24 24"
            >
              <path strokeLinecap="round" strokeLinejoin="round" d="M16.5 4.5l-9 15M4.5 19.5l15-15" />
            </svg>
            Elasticsearch
          </h2>
          <p className="text-gray-700">
            Retrieve logs indexed in Elasticsearch, ideal for full-text search and deeper log analysis.
          </p>
        </Link>

        <div className="bg-white rounded-xl p-6 shadow">
          <h2 className="text-xl font-bold mb-3 flex items-center gap-2">
            <svg
              className="w-6 h-6 text-gray-600"
              fill="none"
              stroke="currentColor"
              strokeWidth="1.5"
              viewBox="0 0 24 24"
            >
              <path strokeLinecap="round" strokeLinejoin="round" d="M7.5 4.5l9 15M4.5 19.5l15-15" />
            </svg>
            PostgreSQL
          </h2>
          <p className="text-gray-700">
            Storing task data in a robust PostgreSQL database, ensuring reliability and ACID compliance.
          </p>
        </div>
      </div>

      <footer className="text-center text-sm text-gray-500 mt-auto py-4">
        © {new Date().getFullYear()} Distributed Task Scheduler. All rights reserved.
      </footer>
    </div>
  );
}
