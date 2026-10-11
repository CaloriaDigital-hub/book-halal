import { http } from "@/shared/api/http";
import type { BookCatalogItem, BookDetails, BookPage, BookPages, UploadBookInput, UploadBookResult } from "../model/types";

const BASE = "/api/backend";
const id = encodeURIComponent;

export async function getCatalog(): Promise<BookCatalogItem[]> {
  return (await http<BookCatalogItem[] | null>(`${BASE}/books/catalog`)) ?? [];
}

export function getBook(bookId: string) {
  return http<BookDetails>(`${BASE}/books/${id(bookId)}`);
}

/** Все страницы книги (limit=0 — без пагинации). URL картинок относительные. */
export async function getAllPages(bookId: string): Promise<BookPage[]> {
  const res = await http<BookPages>(`${BASE}/books/${id(bookId)}/pages?limit=0`);
  return res.pages ?? [];
}

export function uploadBook(input: UploadBookInput) {
  const form = new FormData();
  form.set("title", input.title);
  form.set("author", input.author);
  form.set("description", input.description);
  form.set("price", String(input.price));
  form.set("document", input.document);
  return http<UploadBookResult>(`${BASE}/admin/books/upload`, { method: "POST", body: form });
}
