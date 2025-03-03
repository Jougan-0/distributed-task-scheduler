"use client";
export const dynamic = "force-dynamic";

import { useEffect, useState, useRef } from "react";
import axios from "axios";

interface KafkaEvent {
  event: string;
  task_id?: string;
  task_name?: string;
  task_type?: string;
  priority?: number;
  scheduled_time?: string;
  completed_at?: string;
  failed_at?: string;
  raw?: string;
}
export default function KafkaEventsPage() {
  const [events, setEvents] = useState<KafkaEvent[]>([]);
  const scrollRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const fetchEvents = async () => {
      try {
        const res = await axios.get<KafkaEvent[]>(process.env.NEXT_PUBLIC_BACKEND_URL + "/kafka/events");
        setEvents(res.data);
      } catch (err) {
        console.error("Error fetching Kafka events:", err);
      }
    };

    fetchEvents();
    const interval = setInterval(fetchEvents, 3000);
    return () => clearInterval(interval);
  }, []);

  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
    }
  }, [events]);

  return (
    <div className="p-6 max-w-5xl mx-auto">
      <h1 className="text-3xl font-bold text-white text-center mb-6">Kafka Events</h1>
      <div className="bg-gray-900 p-6 rounded-lg shadow-lg">
        <div ref={scrollRef} className="overflow-y-auto max-h-96 border border-gray-700 rounded-md">
          <table className="w-full border-collapse">
            <thead className="sticky top-0 bg-gray-800">
              <tr>
                <th className="text-left px-4 py-2 text-gray-300">Event</th>
                <th className="text-left px-4 py-2 text-gray-300">Task ID</th>
                <th className="text-left px-4 py-2 text-gray-300">Task Name</th>
                <th className="text-left px-4 py-2 text-gray-300">Type</th>
                <th className="text-left px-4 py-2 text-gray-300">Priority</th>
                <th className="text-left px-4 py-2 text-gray-300">Timestamp</th>
              </tr>
            </thead>
            <tbody>
  {events.length === 0 ? (
    <tr>
      <td colSpan={6} className="text-center text-gray-400 py-4">
        No events yet...
      </td>
    </tr>
  ) : (
    events.map((event, i) => (
      <tr key={`${event.task_id}-${event.event}-${i}`} className={`border-b ${getRowStyle(event.event)}`}>
        <td className="px-4 py-2 font-semibold">{event.event || "Unknown"}</td>
        <td className="px-4 py-2">{event.task_id || "-"}</td>
        <td className="px-4 py-2">{event.task_name || "-"}</td>
        <td className="px-4 py-2">{event.task_type || "-"}</td>
        <td className="px-4 py-2">{event.priority !== undefined ? event.priority : "-"}</td>
        <td className="px-4 py-2 text-gray-400">{formatTimestamp(event)}</td>
      </tr>
    ))
  )}
</tbody>

          </table>
        </div>
      </div>
    </div>
  );
}

function getRowStyle(eventType?: string): string {
  switch (eventType) {
    case "TaskCreated":
      return "bg-blue-700 text-white";
    case "TaskCompleted":
      return "bg-green-700 text-white";
    case "TaskFailed":
      return "bg-red-700 text-white";
    default:
      return "bg-gray-700 text-white";
  }
}

function formatTimestamp(event: KafkaEvent): string {
  return event.completed_at || event.failed_at || event.scheduled_time || "-";
}
