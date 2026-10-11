import { http, jsonInit } from "@/shared/api/http";

export type Role = "admin" | "user";
export type Session = { authenticated: false } | { authenticated: true; role: Role };

export function signIn(email: string, password: string) {
  return http<{ ok: true; role: Role }>("/api/session", jsonInit("POST", { email, password }));
}

export function signOut() {
  return http<null>("/api/session", { method: "DELETE" });
}

export function getSession() {
  return http<Session>("/api/session");
}

export function requestCode(email: string) {
  return http<{ message: string }>("/api/backend/auth/register/initiate", jsonInit("POST", { email }), { authRedirect: false });
}

export type RegisterInput = { email: string; code: string; first_name: string; last_name: string; password: string };

export function confirmRegistration(input: RegisterInput) {
  return http<unknown>("/api/backend/auth/register/confirm", jsonInit("POST", input), { authRedirect: false });
}
