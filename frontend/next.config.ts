import type { NextConfig } from "next";

const backendUrl = process.env.BACKEND_URL ?? "http://localhost:8090";

const nextConfig: NextConfig = {
  cacheComponents: true,
  partialPrefetching: true,
  reactCompiler: true,
  async rewrites() {
    // Обложки и страницы книг отдаёт Go-бэк (публично), проксируем как есть.
    return [{ source: "/static/:path*", destination: `${backendUrl}/static/:path*` }];
  },
};

export default nextConfig;
