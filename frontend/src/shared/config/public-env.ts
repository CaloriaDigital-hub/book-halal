// Публичный адрес Go-бэка. RAG скачивает картинки страниц по абсолютным URL,
// и в rag-service localhost:8090 подменяется на BACKEND_URL контейнера.
export const BACKEND_PUBLIC_URL = process.env.NEXT_PUBLIC_BACKEND_PUBLIC_URL ?? "http://localhost:8090";

export function toBackendPublicUrl(url: string): string {
  return /^https?:\/\//.test(url) ? url : `${BACKEND_PUBLIC_URL}${url.startsWith("/") ? "" : "/"}${url}`;
}
