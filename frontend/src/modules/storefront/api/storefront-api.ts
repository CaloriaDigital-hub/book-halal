import { http } from "@/shared/api/http";
import type { BookDetails, CatalogBook } from "../model/types";

// Каталог и карточка книги публичные — гостя не уводим на вход.
export async function getCatalog(): Promise<CatalogBook[]> {
  return (await http<CatalogBook[] | null>("/api/backend/books/catalog", undefined, { authRedirect: false })) ?? [];
}

export function getBook(bookId: string) {
  return http<BookDetails>(`/api/backend/books/${encodeURIComponent(bookId)}`, undefined, { authRedirect: false });
}
