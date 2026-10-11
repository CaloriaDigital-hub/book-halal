import { http, jsonInit } from "@/shared/api/http";
import type { AskAnswer, Bookmark, ReaderPage } from "../model/types";

const book = (id: string) => `/api/backend/books/${encodeURIComponent(id)}`;

export async function getPages(bookId: string): Promise<ReaderPage[]> {
  const res = await http<{ total_pages: number; pages: ReaderPage[] | null }>(`${book(bookId)}/pages?limit=0`);
  return res.pages ?? [];
}

export async function getProgress(bookId: string): Promise<number> {
  return (await http<{ page_number: number }>(`${book(bookId)}/progress`)).page_number;
}

export function saveProgress(bookId: string, page: number) {
  return http<unknown>(`${book(bookId)}/progress`, jsonInit("PUT", { page_number: page }));
}

export async function getBookmarks(bookId: string): Promise<Bookmark[]> {
  return (await http<Bookmark[] | null>(`${book(bookId)}/bookmarks`)) ?? [];
}

export function addBookmark(bookId: string, page: number) {
  return http<Bookmark>(`${book(bookId)}/bookmarks`, jsonInit("POST", { page_number: page }));
}

export function removeBookmark(bookId: string, page: number) {
  return http<unknown>(`${book(bookId)}/bookmarks/${page}`, { method: "DELETE" });
}

export function ask(bookId: string, question: string) {
  return http<AskAnswer>(`${book(bookId)}/ask`, jsonInit("POST", { question }));
}
