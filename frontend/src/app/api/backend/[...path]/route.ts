import type { NextRequest } from "next/server";
import { getSessionToken } from "@/modules/auth/server";
import { forward } from "@/shared/api/forward";
import { serverEnv } from "@/shared/config/server-env";

// /api/backend/* → Go-бэк /api/*. Токен берём из httpOnly-cookie, права проверяет сам бэк.
async function handler(request: NextRequest, ctx: RouteContext<"/api/backend/[...path]">) {
  const { path } = await ctx.params;
  const token = await getSessionToken();
  if (!token) return Response.json({ error: "Не авторизован" }, { status: 401 });
  return forward(request, `${serverEnv.backendUrl}/api/`, path, { Authorization: `Bearer ${token}` });
}

export { handler as DELETE, handler as GET, handler as PATCH, handler as POST, handler as PUT };
