"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { getPages, getProgress, saveProgress } from "../api/reader-api";
import type { ReaderPage } from "./types";

type State =
  | { status: "loading" }
  | { status: "error"; error: unknown }
  | { status: "ready"; pages: ReaderPage[]; page: number };

const SAVE_DELAY_MS = 800;

/** Страницы книги, текущая страница и сохранение прогресса (с задержкой, чтобы не спамить бэк при листании). */
export function useReader(bookId: string) {
  const [state, setState] = useState<State>({ status: "loading" });
  const saveTimer = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(() => {
    let active = true;
    Promise.all([getPages(bookId), getProgress(bookId).catch(() => 1)]).then(
      ([pages, saved]) => {
        if (active) setState({ status: "ready", pages, page: clamp(saved, pages.length) });
      },
      (error: unknown) => {
        if (active) setState({ status: "error", error });
      },
    );
    return () => {
      active = false;
    };
  }, [bookId]);

  useEffect(() => () => {
    if (saveTimer.current) clearTimeout(saveTimer.current);
  }, []);

  const goTo = useCallback(
    (page: number) => {
      setState((s) => {
        if (s.status !== "ready") return s;
        const next = clamp(page, s.pages.length);
        if (next === s.page) return s;
        if (saveTimer.current) clearTimeout(saveTimer.current);
        saveTimer.current = setTimeout(() => saveProgress(bookId, next).catch(() => {}), SAVE_DELAY_MS);
        return { ...s, page: next };
      });
    },
    [bookId],
  );

  return { state, goTo };
}

function clamp(page: number, total: number) {
  return Math.min(Math.max(1, page || 1), Math.max(1, total));
}
