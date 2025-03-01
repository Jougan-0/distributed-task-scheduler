"use client";
import { useEffect, useState } from "react";
import axios from "axios";

type ElasticHit = {
  _index: string;
  _id: string;
  _source: any;
};

export default function ElasticsearchPage() {
  const [logs, setLogs] = useState<ElasticHit[]>([]);
  const ELASTICSEARCH_URL = process.env.NEXT_PUBLIC_ELASTICSEARCH_URL;
  useEffect(() => {
    const fetchLogs = async () => {
        const res = await axios.get(`${ELASTICSEARCH_URL}/logs/_search`);
        setLogs(res.data.hits?.hits || []);
    };
    fetchLogs();
  }, []);

  return (
    <div>
      <h1 className="text-2xl font-bold">Elasticsearch Logs</h1>
      <div className="bg-gray-800 text-white p-4 mt-4 rounded-md">
        {logs.length ? (
          logs.map((log, i) => (
            <p key={i} className="mb-2">
              {JSON.stringify(log)}
            </p>
          ))
        ) : (
          <p>No Elasticsearch logs found...</p>
        )}
      </div>
    </div>
  );
}
