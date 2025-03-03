"use client";
import { useEffect, useState } from "react";
import axios from "axios";

interface RedisKeyValue {
  key: string;
  value: string;
}
export default function RedisPage() {
  const [data, setData] = useState<RedisKeyValue[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const BACKEND_URL = process.env.NEXT_PUBLIC_BACKEND_URL || "http://localhost:8080";

  useEffect(() => {
    let intervalId: NodeJS.Timeout;

    const fetchData = async () => {
      try {
        const res = await axios.get(`${BACKEND_URL}/redis/keys`);
        const responseData = res.data;
        console.log(res)
        if (res.status==200){
          if (res.data==null){
            setData([]);
          } else {
            const validData = responseData.filter((item: any) => item?.key && item?.value);
            setData(validData);
          }
        } else {
          throw new Error(`Request failed with status ${res.status}`);
        }
      
        setError(null);
      } catch (err) {
        setError("Failed to fetch Redis data. Please try again.");
        console.error("Error fetching Redis keys:", err);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
    intervalId = setInterval(fetchData, 5000);
    return () => clearInterval(intervalId);
  }, [BACKEND_URL]);

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-4 text-center text-white">Redis Key-Value Store</h1>

      {loading ? (
        <p className="text-gray-400 text-center">Loading Redis data...</p>
      ) : error ? (
        <p className="text-red-500 text-center">{error}</p>
      ) : data.length === 0 ? (
        <p className="text-gray-400 text-center">Cache is cleared. Nothing is being processed.</p>
      ) : (
        <div className="overflow-x-auto bg-gray-900 text-white p-4 rounded-lg shadow-md">
          <table className="w-full border-collapse">
            <thead>
              <tr className="bg-gray-700 text-left">
                <th className="p-3 border border-gray-600">Key</th>
                <th className="p-3 border border-gray-600">Value</th>
              </tr>
            </thead>
            <tbody>
              {data.map((item, index) => (
                <tr key={index} className="border border-gray-700 hover:bg-gray-800">
                  <td className="p-3">{item.key || "N/A"}</td>
                  <td className="p-3">{item.value || "N/A"}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
