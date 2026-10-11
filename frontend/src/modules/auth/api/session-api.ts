import { http, jsonInit } from "@/shared/api/http";

export function signIn(email: string, password: string) {
  return http<{ ok: true }>("/api/session", jsonInit("POST", { email, password }));
}

export function signOut() {
  return http<null>("/api/session", { method: "DELETE" });
}
