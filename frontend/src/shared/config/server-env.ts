import "server-only";

export const serverEnv = {
  backendUrl: process.env.BACKEND_URL ?? "http://localhost:8090",
  ragUrl: process.env.RAG_URL ?? "http://localhost:8000",
};
