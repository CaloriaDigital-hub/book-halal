import type { NextRequest } from "next/server";
import { getSessionToken } from "@/modules/auth/server";
import { forward } from "@/shared/api/forward";
import { serverEnv } from "@/shared/config/server-env";

// /api/backend/* → Go-бэк /api/*. Токен берём из httpOnly-cookie; без сессии запрос идёт
// анонимно — публичные эндпоинты ответят, закрытые вернут 401. Права проверяет сам бэк.
async function handler(request: NextRequest, ctx: RouteContext<"/api/backend/[...path]">) {
  const { path } = await ctx.params;
  const token = await getSessionToken();
  return forward(request, `${serverEnv.backendUrl}/api/`, path, token ? { Authorization: `Bearer ${token}` } : {});
}

export { handler as DELETE, handler as GET, handler as PATCH, handler as POST, handler as PUT };
