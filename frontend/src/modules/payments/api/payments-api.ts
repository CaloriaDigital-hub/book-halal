import { http } from "@/shared/api/http";
import type { AdminPurchase, PaymentSettings, PurchaseStatus } from "../model/types";

const BASE = "/api/backend";
const id = encodeURIComponent;

export function getRequisites() {
  return http<PaymentSettings>(`${BASE}/payments/requisites`);
}

export function updateRequisites(input: { details: string; qr?: File | null; removeQr?: boolean }) {
  const form = new FormData();
  form.set("details", input.details);
  if (input.qr) form.set("qr", input.qr);
  if (input.removeQr) form.set("remove_qr", "true");
  return http<PaymentSettings>(`${BASE}/admin/payments/requisites`, { method: "PUT", body: form });
}

export async function listPurchases(status: PurchaseStatus | "all") {
  const query = status === "all" ? "" : `?status=${status}`;
  return (await http<AdminPurchase[] | null>(`${BASE}/admin/purchases${query}`)) ?? [];
}

export function reviewPurchase(purchaseId: string, approve: boolean) {
  return http<AdminPurchase>(`${BASE}/admin/purchases/${id(purchaseId)}/${approve ? "approve" : "reject"}`, {
    method: "POST",
  });
}

/** Чек отдаётся только админам — через BFF с cookie сессии. */
export function receiptUrl(purchaseId: string) {
  return `${BASE}/admin/purchases/${id(purchaseId)}/receipt`;
}
