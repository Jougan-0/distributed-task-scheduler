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

  const [backendUrl, setBackendUrl] = useState("");

  useEffect(() => {
    fetch("/api/config")
      .then((res) => res.json())
      .then((data) => {
        setBackendUrl(data.backendUrl);
        console.log("Backend URL set to:", data.backendUrl);
      })
      .catch((err) => console.error("Error fetching config:", err));
  }, []);

  const createExampleTask = async () => {
    if (!backendUrl) {
      setMessage("Backend URL not set yet.");
      return;
    }

    setMessage("Creating example task...");
    try {
      const res = await fetch(`${backendUrl}/api/v1/tasks`, {
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
    if (!backendUrl) {
      setMessage("Backend URL not set yet.");
      return;
    }

    setMessage("Creating custom task...");
    try {
      const res = await fetch(`${backendUrl}/api/v1/tasks`, {
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
        const errorData = await res.json();
        throw new Error(`Failed to create custom task: ${errorData.message || res.statusText}`);
      }

      setMessage(`Custom Task created successfully! (Type: ${taskType}, Email: ${taskEmail}, Priority: ${priority}). Check logs/metrics pages.`);
    } catch (err: any) {
      console.error("Error:", err);
      setMessage("Error creating custom task: " + (err.message || err));
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

      <footer className="text-center text-sm text-gray-500 mt-auto py-4">
        © {new Date().getFullYear()} Distributed Task Scheduler. All rights reserved.
      </footer>
    </div>
  );
}
