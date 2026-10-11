import type { NextRequest } from "next/server";
import { getSessionToken, requireAdmin } from "@/modules/auth/server";
import { forward } from "@/shared/api/forward";
import { serverEnv } from "@/shared/config/server-env";

// /api/rag/* → RAG-сервис. У самого RAG авторизации нет, поэтому пускаем только админов.
async function handler(request: NextRequest, ctx: RouteContext<"/api/rag/[...path]">) {
  const { path } = await ctx.params;
  const access = await requireAdmin(await getSessionToken());
  if (access instanceof Response) return access;
  return forward(request, serverEnv.ragUrl, path);
}

export { handler as DELETE, handler as GET, handler as POST };
