import { NextResponse, type NextRequest } from "next/server";
import { checkAdmin, SESSION_COOKIE } from "@/modules/auth/server";

// Без сессии — на вход (и обратно после него). В /admin пускаем только админов.
// Реальные права всё равно проверяют BFF и Go-бэк на каждом запросе.
export async function proxy(request: NextRequest) {
  const token = request.cookies.get(SESSION_COOKIE)?.value;
  if (!token) {
    const url = new URL("/login", request.url);
    url.searchParams.set("next", request.nextUrl.pathname + request.nextUrl.search);
    return NextResponse.redirect(url);
  }
  if (request.nextUrl.pathname.startsWith("/admin") && (await checkAdmin(token)) === "forbidden") {
    return NextResponse.redirect(new URL("/", request.url));
  }
  return NextResponse.next();
}

export const config = {
  matcher: ["/admin/:path*", "/books/:id/read", "/books/:id/buy"],
};
