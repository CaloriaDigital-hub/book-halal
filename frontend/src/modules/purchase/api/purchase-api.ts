import { ApiError, http } from "@/shared/api/http";
import type { ContactType, MyPurchase, Requisites } from "../model/types";

const book = (id: string) => `/api/backend/books/${encodeURIComponent(id)}`;

export function getRequisites() {
  return http<Requisites>("/api/backend/payments/requisites");
}

/** Последняя оплата книги текущим пользователем или null, если ещё не платил. */
export async function getMyPurchase(bookId: string): Promise<MyPurchase | null> {
  try {
    return await http<MyPurchase>(`${book(bookId)}/purchase`);
  } catch (e) {
    if (e instanceof ApiError && e.status === 404) return null;
    throw e;
  }
}

export function submitPurchase(bookId: string, input: { receipt: File; contactType: ContactType; contact: string }) {
  const form = new FormData();
  form.set("receipt", input.receipt);
  form.set("contact_type", input.contactType);
  form.set("contact", input.contact);
  return http<MyPurchase>(`${book(bookId)}/purchase`, { method: "POST", body: form });
}
