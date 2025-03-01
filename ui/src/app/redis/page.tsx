"use client";
import { useEffect, useState } from "react";
import axios from "axios";

export default function RedisPage() {
  const [keys, setKeys] = useState<string[]>([]);
  const BACKEND_URL = process.env.NEXT_PUBLIC_BACKEND_URL || "http://localhost:8080";

  useEffect(() => {
    const fetchKeys = async () => {
      try {
        const res = await axios.get(`${BACKEND_URL}/redis/keys`);
        setKeys(res.data);
      } catch (err) {
        console.error("Error fetching Redis keys:", err);
      }
    };

    fetchKeys();
  }, [BACKEND_URL]);

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-4">Redis Keys</h1>
      <div className="bg-gray-800 text-white p-4 mt-4 rounded-md">
        {keys.length === 0 ? (
          <p className="text-gray-400">No keys found in Redis.</p>
        ) : (
          keys.map((key, i) => <p key={i}>{key}</p>)
        )}
      </div>
    </div>
  );
}
