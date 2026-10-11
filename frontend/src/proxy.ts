import { NextResponse, type NextRequest } from "next/server";
import { SESSION_COOKIE } from "@/modules/auth/server";

// Оптимистичная проверка: нет cookie — сразу на логин. Реальные права
// проверяются в BFF-роутах (/api/backend, /api/rag) на каждом запросе.
export function proxy(request: NextRequest) {
  if (request.cookies.has(SESSION_COOKIE)) return NextResponse.next();
  const url = new URL("/login", request.url);
  url.searchParams.set("next", request.nextUrl.pathname + request.nextUrl.search);
  return NextResponse.redirect(url);
}

export const config = {
  matcher: "/admin/:path*",
};
