import { http, jsonInit } from "@/shared/api/http";
import type { AskResult, IndexPageInput, IndexResult, IndexSnapshot, PageChunk, RagHealth } from "../model/types";

const BASE = "/api/rag";
const id = encodeURIComponent;

export function getHealth() {
  return http<RagHealth>(`${BASE}/health`);
}

export function inspectIndex(bookId: string, limit: number) {
  return http<IndexSnapshot>(`${BASE}/admin/books/${id(bookId)}/index?limit=${limit}`);
}

export function deleteIndex(bookId: string) {
  return http<{ deleted: string }>(`${BASE}/admin/books/${id(bookId)}/index`, { method: "DELETE" });
}

export function indexPages(bookId: string, pages: IndexPageInput[]) {
  return http<IndexResult>(`${BASE}/admin/books/${id(bookId)}/index`, jsonInit("POST", { pages }));
}

export function askBook(bookId: string, question: string) {
  return http<AskResult>(`${BASE}/books/${id(bookId)}/ask`, jsonInit("POST", { question }));
}

export function getPageChunk(bookId: string, page: number) {
  return http<PageChunk>(`${BASE}/debug/page/${id(bookId)}?page=${page}`);
}
