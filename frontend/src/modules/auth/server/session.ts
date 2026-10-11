import "server-only";
import { cookies } from "next/headers";

// Одна сессия на покупателя и админа; роль проверяется в BFF.
export const SESSION_COOKIE = "bh_session";

export async function getSessionToken(): Promise<string | null> {
  return (await cookies()).get(SESSION_COOKIE)?.value ?? null;
}

export async function setSessionToken(token: string, expiresAt?: string) {
  const expires = expiresAt ? new Date(expiresAt) : undefined;
  (await cookies()).set(SESSION_COOKIE, token, {
    httpOnly: true,
    sameSite: "lax",
    secure: process.env.NODE_ENV === "production",
    path: "/",
    expires: expires && !Number.isNaN(expires.getTime()) ? expires : undefined,
  });
}

export async function clearSessionToken() {
  (await cookies()).delete(SESSION_COOKIE);
}
