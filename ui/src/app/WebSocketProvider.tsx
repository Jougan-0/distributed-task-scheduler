"use client";
import React, { createContext, useContext, useState, useEffect } from "react";

interface IWebSocketContext {
  logs: string[];
}

const WebSocketContext = createContext<IWebSocketContext | null>(null);

export function WebSocketProvider({ children }: { children: React.ReactNode }) {
  const [logs, setLogs] = useState<string[]>([]);

  useEffect(() => {
    const wsUrl = process.env.NEXT_PUBLIC_BACKEND_URL
      ? process.env.NEXT_PUBLIC_BACKEND_URL.replace(/^http/, "ws") + "/ws"
      : "ws://localhost:8080/ws";

    const ws = new WebSocket(wsUrl);

    ws.onopen = () => {
      console.log("WebSocket connected (global)");
    };

    ws.onmessage = (event) => {
      setLogs((prev) => [event.data, ...prev]);
    };

    ws.onerror = (err) => {
      console.error("WebSocket error:", err);
    };

    ws.onclose = () => {
      console.log("WebSocket closed");
    };

    return () => {
      ws.close();
    };
  }, []);

  return (
    <WebSocketContext.Provider value={{ logs }}>
      {children}
    </WebSocketContext.Provider>
  );
}

export function useWebSocketContext() {
  const ctx = useContext(WebSocketContext);
  if (!ctx) {
    throw new Error("useWebSocketContext must be used within a WebSocketProvider");
  }
  return ctx;
}
