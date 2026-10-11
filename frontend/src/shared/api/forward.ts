import "server-only";
import type { NextRequest } from "next/server";
import { Agent, fetch as undiciFetch, type RequestInit as UndiciRequestInit } from "undici";

const FORWARDED_REQUEST_HEADERS = ["content-type", "accept", "accept-language"];
// Тело ответа уже раскодировано fetch'ем, эти заголовки стали бы враньём.
const DROPPED_RESPONSE_HEADERS = ["content-encoding", "content-length", "transfer-encoding", "connection"];

// У встроенного fetch в Node таймаут на заголовки ответа 5 минут, а индексация книги
// (OCR всех страниц в одном запросе) идёт дольше. Даём апстриму до 30 минут.
const UPSTREAM_TIMEOUT_MS = 30 * 60_000;
const upstreamAgent = new Agent({ headersTimeout: UPSTREAM_TIMEOUT_MS, bodyTimeout: UPSTREAM_TIMEOUT_MS });

/**
 * Проксирует запрос в апстрим как есть (метод, query, тело-стрим),
 * добавляя только переданные заголовки (например, Authorization).
 */
export async function forward(
  req: NextRequest,
  base: string,
  path: string[],
  extraHeaders: Record<string, string> = {},
): Promise<Response> {
  const target = new URL(path.map(encodeURIComponent).join("/"), withTrailingSlash(base));
  target.search = req.nextUrl.search;

  const headers = new Headers(extraHeaders);
  for (const name of FORWARDED_REQUEST_HEADERS) {
    const value = req.headers.get(name);
    if (value) headers.set(name, value);
  }

  const init: UndiciRequestInit = {
    method: req.method,
    headers,
    redirect: "manual",
    dispatcher: upstreamAgent,
  };
  if (req.method !== "GET" && req.method !== "HEAD" && req.body) {
    init.body = req.body as unknown as UndiciRequestInit["body"];
    init.duplex = "half";
  }

  let upstream: Awaited<ReturnType<typeof undiciFetch>>;
  try {
    upstream = await undiciFetch(target, init);
  } catch {
    return Response.json({ error: `Сервис недоступен: ${target.origin}` }, { status: 502 });
  }

  const resHeaders = new Headers(upstream.headers as unknown as HeadersInit);
  for (const name of DROPPED_RESPONSE_HEADERS) resHeaders.delete(name);
  return new Response(upstream.body as ReadableStream | null, { status: upstream.status, headers: resHeaders });
}

function withTrailingSlash(url: string): string {
  return url.endsWith("/") ? url : `${url}/`;
}
