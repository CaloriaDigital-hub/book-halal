import { checkAdmin, clearSessionToken, getSessionToken, setSessionToken } from "@/modules/auth/server";
import { serverEnv } from "@/shared/config/server-env";

type SignInResponse = { token: string; expires_at?: string };
type Role = "admin" | "user";

/** Вход покупателя или админа. Роль возвращаем, чтобы фронт знал, куда вести после входа. */
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
    return Response.json({ error: "Сервис недоступен, попробуйте позже" }, { status: 502 });
  }
  if (res.status === 401) return Response.json({ error: "Неверный email или пароль" }, { status: 401 });
  if (!res.ok) return Response.json({ error: "Не удалось войти" }, { status: res.status });

  const session = (await res.json()) as SignInResponse;
  await setSessionToken(session.token, session.expires_at);
  return Response.json({ ok: true, role: await roleOf(session.token) });
}

/** Текущая сессия: гость, покупатель или админ. */
export async function GET() {
  const token = await getSessionToken();
  if (!token) return Response.json({ authenticated: false });
  const check = await checkAdmin(token);
  if (check === "unauthenticated") {
    await clearSessionToken();
    return Response.json({ authenticated: false });
  }
  return Response.json({ authenticated: true, role: check === "admin" ? "admin" : "user" });
}

export async function DELETE() {
  await clearSessionToken();
  return new Response(null, { status: 204 });
}

async function roleOf(token: string): Promise<Role> {
  return (await checkAdmin(token)) === "admin" ? "admin" : "user";
}
