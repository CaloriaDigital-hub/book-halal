"use client";

import { useAsync } from "@/shared/lib/use-async";
import { getBook } from "../api/storefront-api";

/** Карточка книги для экранов, которые собирают несколько модулей (оплата, читалка). */
export function useStoreBook(bookId: string) {
  return useAsync(`store-book:${bookId}`, () => getBook(bookId));
}
