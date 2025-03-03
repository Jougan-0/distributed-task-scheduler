import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  reactStrictMode: true, 
  experimental: {
    serverActions: {}, 
  },
};

export default nextConfig;
