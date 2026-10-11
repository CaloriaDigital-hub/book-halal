// Клиентский HTTP-хелпер. Все запросы идут в BFF этого же Next-приложения
// (/api/backend/*, /api/rag/*, /api/session), токен живёт в httpOnly-cookie.

export class ApiError extends Error {
  constructor(
    readonly status: number,
    message: string,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

const LOGIN_PATH = "/login";

type HttpOptions = {
  /** false — не уводить на /login при 401 (публичные страницы, где гость — норма). */
  authRedirect?: boolean;
};

export async function http<T>(url: string, init?: RequestInit, options: HttpOptions = {}): Promise<T> {
  let res: Response;
  try {
    res = await fetch(url, { cache: "no-store", ...init });
  } catch {
    throw new ApiError(0, "Сеть недоступна");
  }

  const text = await res.text();
  const data = text ? parseJson(text) : null;

  if (res.status === 401 && options.authRedirect !== false && !url.startsWith("/api/session")) {
    redirectToLogin();
  }
  if (!res.ok) {
    throw new ApiError(res.status, errorMessage(data, text, res.status));
  }
  return data as T;
}

export function jsonInit(method: string, body?: unknown): RequestInit {
  return {
    method,
    headers: { "Content-Type": "application/json" },
    body: body === undefined ? undefined : JSON.stringify(body),
  };
}

function parseJson(text: string): unknown {
  try {
    return JSON.parse(text);
  } catch {
    return text;
  }
}

// Go-бэк отвечает {"error": "..."}, FastAPI — {"detail": "..." | [...]}.
function errorMessage(data: unknown, text: string, status: number): string {
  if (data && typeof data === "object") {
    const { error, detail } = data as { error?: unknown; detail?: unknown };
    if (typeof error === "string") return error;
    if (typeof detail === "string") return detail;
    if (detail !== undefined) return JSON.stringify(detail);
  }
  return text.trim() || `HTTP ${status}`;
}

function redirectToLogin() {
  if (typeof window === "undefined" || window.location.pathname === LOGIN_PATH) return;
  const next = encodeURIComponent(window.location.pathname + window.location.search);
  // Полная перезагрузка намеренно: сбрасываем всё клиентское состояние протухшей сессии.
  // eslint-disable-next-line @next/next/no-location-assign-relative-destination
  window.location.assign(`${LOGIN_PATH}?next=${next}`);
}
