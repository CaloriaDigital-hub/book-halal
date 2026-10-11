"use client";

import { useState } from "react";
import { deleteIndex, indexPages } from "../api/rag-api";
import type { IndexPageInput } from "./types";

export type ReindexPhase = "idle" | "loading-pages" | "deleting" | "indexing" | "done" | "error";

export type ReindexState = {
  phase: ReindexPhase;
  total: number;
  indexed: number;
  skipped: number;
  chunks?: number;
  startedAt?: number;
  error?: string;
};

const IDLE: ReindexState = { phase: "idle", total: 0, indexed: 0, skipped: 0 };

/**
 * Полная переиндексация: удалить коллекцию → отправить ВСЕ страницы одним запросом.
 * Пачками слать нельзя: RAG на каждый вызов index() перестраивает граф книги заново
 * и нумерует точки с нуля, так что каждая пачка затёрла бы предыдущую.
 */
export function useReindex(bookId: string, loadPages: () => Promise<IndexPageInput[]>, onChanged?: () => void) {
  const [state, setState] = useState<ReindexState>(IDLE);
  const running = state.phase === "loading-pages" || state.phase === "deleting" || state.phase === "indexing";

  async function start() {
    setState({ ...IDLE, phase: "loading-pages" });
    let touched = false;
    try {
      const pages = await loadPages();
      if (pages.length === 0) throw new Error("У книги нет страниц — индексировать нечего");

      setState((s) => ({ ...s, phase: "deleting", total: pages.length }));
      await deleteIndex(bookId);
      touched = true;

      setState((s) => ({ ...s, phase: "indexing", startedAt: Date.now() }));
      const res = await indexPages(bookId, pages);
      setState((s) => ({ ...s, phase: "done", indexed: res.indexed, skipped: res.skipped, chunks: res.chunks }));
    } catch (e) {
      setState((s) => ({ ...s, phase: "error", error: e instanceof Error ? e.message : String(e) }));
    } finally {
      if (touched) onChanged?.();
    }
  }

  return { state, running, start };
}
