"use client";
export const dynamic = "force-dynamic";

import { useWebSocketContext } from "../WebSocketProvider";
export default function LogsPage() {
  const { logs } = useWebSocketContext();

  return (
    <div className="min-h-screen bg-gray-900 text-white p-4">
      <h1 className="text-2xl font-bold text-white">Backend Logs (WebSocket)</h1>
      <div className="bg-gray-800 p-4 mt-4 rounded-lg min-h-[300px]">
        {logs.length === 0 ? (
          <p className="text-gray-400">No logs yet...</p>
        ) : (
          <ul className="space-y-2">
            {logs.map((log, index) => (
              <li key={index} className="p-2 bg-gray-700 rounded">
                {log}
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}
