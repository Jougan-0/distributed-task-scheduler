"use client";

import { useEffect, useState } from "react";
import axios from "axios";
import { Loader2 } from "lucide-react";

interface Task {
  Attempts: number;
  CreatedAt: string;
  ID: string;
  MaxRetries: number;
  Name: string;
  Payload: string;
  Priority: number;
  ScheduledTime: string;
  Status: string;
  Type: string;
  UpdatedAt: string;
}

export default function ElasticsearchPage() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [query, setQuery] = useState("DemoTask");
  const [loading, setLoading] = useState(false);

  const BACKEND_URL = process.env.NEXT_PUBLIC_BACKEND_URL || "http://localhost:8080";

  const fetchTasks = async () => {
    setLoading(true);
    try {
      const res = await axios.get(`${BACKEND_URL}/api/v1/tasks/search/${query}`);
      setTasks(res.data || []);
    } catch (error) {
      console.error("Error fetching tasks:", error);
      setTasks([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchTasks();
  }, []);

  return (
    <div className="max-w-6xl mx-auto p-6">
      <h1 className="text-3xl font-bold mb-6 text-center text-white">Elasticsearch Tasks</h1>

      <div className="flex gap-4 mb-6">
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          className="flex-grow p-3 border rounded-md bg-gray-800 text-white focus:outline-none focus:ring focus:ring-blue-500"
          placeholder="Enter query (e.g., DemoTask)"
        />
        <button
          onClick={fetchTasks}
          className="px-5 py-3 bg-blue-600 text-white rounded-md hover:bg-blue-700 flex items-center"
        >
          {loading ? <Loader2 className="animate-spin w-5 h-5" /> : "Search"}
        </button>
      </div>

      <div className="bg-gray-900 text-white p-5 rounded-lg shadow-md">
        {loading ? (
          <div className="flex justify-center py-6">
            <Loader2 className="animate-spin w-8 h-8 text-blue-500" />
          </div>
        ) : tasks.length > 0 ? (
          <div className="overflow-x-auto">
            <table className="w-full text-left border-collapse">
              <thead>
                <tr className="bg-gray-800 text-gray-300">
                  <th className="p-3 border-b border-gray-700">ID</th>
                  <th className="p-3 border-b border-gray-700">Created</th>
                  <th className="p-3 border-b border-gray-700">Updated</th>
                  <th className="p-3 border-b border-gray-700">Name</th>
                  <th className="p-3 border-b border-gray-700">Payload</th>
                  <th className="p-3 border-b border-gray-700">Status</th>
                  <th className="p-3 border-b border-gray-700">Attempts</th>
                  <th className="p-3 border-b border-gray-700">MaxRetries</th>
                  <th className="p-3 border-b border-gray-700">Priority</th>
                  <th className="p-3 border-b border-gray-700">Scheduled</th>
                  <th className="p-3 border-b border-gray-700">Type</th>
                </tr>
              </thead>
              <tbody>
                {tasks.map((task, index) => (
                  <tr key={index} className="hover:bg-gray-800 border-b border-gray-700">
                    <td className="p-3">{task.ID}</td>
                    <td className="p-3">{new Date(task.CreatedAt).toLocaleString()}</td>
                    <td className="p-3">{new Date(task.UpdatedAt).toLocaleString()}</td>
                    <td className="p-3">{task.Name}</td>
                    <td className="p-3 whitespace-pre-wrap max-w-xs truncate">{task.Payload}</td>
                    <td className="p-3">{task.Status}</td>
                    <td className="p-3">{task.Attempts}</td>
                    <td className="p-3">{task.MaxRetries}</td>
                    <td className="p-3">{task.Priority}</td>
                    <td className="p-3">{new Date(task.ScheduledTime).toLocaleString()}</td>
                    <td className="p-3">{task.Type}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : (
          <p className="text-center py-4">No tasks found.</p>
        )}
      </div>
    </div>
  );
}
