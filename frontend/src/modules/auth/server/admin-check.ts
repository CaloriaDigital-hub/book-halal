import "server-only";
import { serverEnv } from "@/shared/config/server-env";

export type AdminCheck = "admin" | "forbidden" | "unauthenticated" | "unavailable";

const TTL_MS = 60_000;
const cache = new Map<string, { result: AdminCheck; until: number }>();

/**
 * В Go-бэке нет /me, а sign-in не отдаёт роль. Поэтому роль проверяем пробой
 * admin-эндпоинта с пустым телом: RequireAdmin отвечает 401/403 ещё до разбора
 * формы, а админу бэк вернёт 400 «failed to parse form data». Ничего не создаётся.
 * TODO: заменить на GET /api/auth/me, когда он появится в бэке.
 */
export async function checkAdmin(token: string): Promise<AdminCheck> {
  const hit = cache.get(token);
  if (hit && hit.until > Date.now()) return hit.result;

  let status: number;
  try {
    const res = await fetch(`${serverEnv.backendUrl}/api/admin/books/upload`, {
      method: "POST",
      headers: { Authorization: `Bearer ${token}` },
      cache: "no-store",
    });
    status = res.status;
  } catch {
    return "unavailable";
  }

  const result: AdminCheck =
    status === 400 ? "admin" : status === 403 ? "forbidden" : status === 401 ? "unauthenticated" : "unavailable";
  if (result === "admin" || result === "forbidden") {
    cache.set(token, { result, until: Date.now() + TTL_MS });
  }
  return result;
}

const DENIED: Record<Exclude<AdminCheck, "admin">, { status: number; error: string }> = {
  unauthenticated: { status: 401, error: "Сессия истекла, войдите заново" },
  forbidden: { status: 403, error: "Нужна роль admin" },
  unavailable: { status: 502, error: "Бэкенд недоступен — не удалось проверить права" },
};

/** Для route handlers: токен админа или готовый ответ с ошибкой. */
export async function requireAdmin(token: string | null): Promise<{ token: string } | Response> {
  if (!token) return Response.json({ error: "Не авторизован" }, { status: 401 });
  const result = await checkAdmin(token);
  if (result === "admin") return { token };
  return Response.json({ error: DENIED[result].error }, { status: DENIED[result].status });
}
