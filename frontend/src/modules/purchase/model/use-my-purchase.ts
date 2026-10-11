"use client";

import { useEffect } from "react";
import { useAsync } from "@/shared/lib/use-async";
import { getMyPurchase } from "../api/purchase-api";

const POLL_MS = 15_000;

/** Статус оплаты; пока чек на проверке — опрашиваем бэк, чтобы сразу показать результат. */
export function useMyPurchase(bookId: string, enabled: boolean) {
  const result = useAsync(`my-purchase:${bookId}:${enabled}`, () => (enabled ? getMyPurchase(bookId) : Promise.resolve(null)));
  const pending = result.data?.status === "pending";
  const { reload } = result;

  useEffect(() => {
    if (!pending) return;
    const timer = setInterval(reload, POLL_MS);
    return () => clearInterval(timer);
  }, [pending, reload]);

  return result;
}
