import { checkAdmin, clearSessionToken, setSessionToken } from "@/modules/auth/server";
import { serverEnv } from "@/shared/config/server-env";

type SignInResponse = { token: string; expires_at?: string };

export async function POST(request: Request) {
  const body = (await request.json().catch(() => null)) as { email?: unknown; password?: unknown } | null;
  if (typeof body?.email !== "string" || typeof body.password !== "string") {
    return Response.json({ error: "Нужны email и пароль" }, { status: 400 });
  }

  let res: Response;
  try {
    res = await fetch(`${serverEnv.backendUrl}/api/auth/sign-in`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email: body.email, password: body.password }),
      cache: "no-store",
    });
  } catch {
    return Response.json({ error: "Бэкенд недоступен" }, { status: 502 });
  }
  if (res.status === 401) return Response.json({ error: "Неверный email или пароль" }, { status: 401 });
  if (!res.ok) return Response.json({ error: "Не удалось войти" }, { status: res.status });

  const session = (await res.json()) as SignInResponse;
  const role = await checkAdmin(session.token);
  if (role === "forbidden") return Response.json({ error: "У этого аккаунта нет роли admin" }, { status: 403 });
  if (role !== "admin") return Response.json({ error: "Не удалось проверить права" }, { status: 502 });

  await setSessionToken(session.token, session.expires_at);
  return Response.json({ ok: true });
}

export async function DELETE() {
  await clearSessionToken();
  return new Response(null, { status: 204 });
}
