"use client";

import { Alert, EmptyState, Spinner } from "@/shared/ui";
import type { Bookmark } from "../model/types";
import styles from "./BookmarksPanel.module.css";

type Props = {
  bookmarks: Bookmark[] | undefined;
  error?: Error;
  currentPage: number;
  onOpenPage: (page: number) => void;
};

export function BookmarksPanel({ bookmarks, error, currentPage, onOpenPage }: Props) {
  if (error) return <Alert tone="danger">{error.message}</Alert>;
  if (!bookmarks) return <Spinner label="Загрузка…" />;
  if (bookmarks.length === 0) return <EmptyState>Закладок пока нет. Нажмите 🔖 на нужной странице.</EmptyState>;

  return (
    <ul className={styles.list}>
      {[...bookmarks]
        .sort((a, b) => a.page_number - b.page_number)
        .map((b) => (
          <li key={b.id}>
            <button
              type="button"
              className={b.page_number === currentPage ? `${styles.item} ${styles.current}` : styles.item}
              onClick={() => onOpenPage(b.page_number)}
            >
              <span>🔖 Страница {b.page_number}</span>
              <span className={styles.date}>{new Date(b.created_at).toLocaleDateString("ru-RU")}</span>
            </button>
          </li>
        ))}
    </ul>
  );
}
