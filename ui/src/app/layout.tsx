import Link from "next/link";
import "./globals.css";
import { WebSocketProvider } from "./WebSocketProvider";

export const metadata = {
  title: "Distributed Task Scheduler",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body>
        <WebSocketProvider>
          <nav className="bg-gray-900 text-white p-4">
            <div className="container mx-auto flex justify-around">
              <Link href="/">Home</Link>
              <Link href="/logs">Logs</Link>
              <Link href="/metrics">Metrics</Link>
              <Link href="/kafka">Kafka</Link>
              <Link href="/redis">Redis</Link>
              <Link href="/elasticsearch">Elasticsearch</Link>
            </div>
          </nav>
          <main className="p-6">{children}</main>
        </WebSocketProvider>
        <footer className="text-center text-m text-gray-200 mt-auto">
      © {new Date().getFullYear()} Made with ❤️ by Jougan.
            </footer>
      </body>
    </html>
  );
}
