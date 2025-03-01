"use client";
import { useEffect, useState } from "react";
import axios from "axios";

export default function KafkaEventsPage() {
  const [events, setEvents] = useState<string[]>([]);

  useEffect(() => {
    const fetchEvents = async () => {
      try {
        const res = await axios.get(process.env.NEXT_PUBLIC_BACKEND_URL + "/kafka/events");
        setEvents(res.data);
      } catch (err) {
        console.error("Error fetching Kafka events:", err);
      }
    };

    fetchEvents();
    const interval = setInterval(fetchEvents, 3000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold">Kafka Events</h1>
      <div className="bg-gray-800 text-white p-4 mt-4 rounded-md">
        {events.length === 0 ? (
          <p className="text-gray-400">No events yet...</p>
        ) : (
          events.map((ev, i) => <div key={i}>{ev}</div>)
        )}
      </div>
    </div>
  );
}
