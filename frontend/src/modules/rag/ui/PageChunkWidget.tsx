"use client";

import { useState, type FormEvent } from "react";
import { Alert, Button, Input, Widget } from "@/shared/ui";
import { getPageChunk } from "../api/rag-api";
import type { PageChunk } from "../model/types";
import styles from "./PageChunkWidget.module.css";

/** Что именно OCR достал со страницы — главный инструмент при плохих ответах. */
export function PageChunkWidget({ bookId }: { bookId: string }) {
  const [page, setPage] = useState("1");
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [chunk, setChunk] = useState<(PageChunk & { requested: number }) | null>(null);

  async function onSubmit(event: FormEvent) {
    event.preventDefault();
    const n = Number(page);
    if (!Number.isInteger(n) || n < 1) return;
    setPending(true);
    setError(null);
    try {
      setChunk({ ...(await getPageChunk(bookId, n)), requested: n });
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    } finally {
      setPending(false);
    }
  }

  return (
    <Widget title="Текст страницы в индексе" description="Проверка качества OCR по конкретной странице">
      <form className={styles.form} onSubmit={onSubmit}>
        <Input type="number" min={1} value={page} onChange={(e) => setPage(e.target.value)} aria-label="Номер страницы" />
        <Button type="submit" loading={pending}>
          Показать
        </Button>
      </form>
      {error && <Alert tone="danger">{error}</Alert>}
      {chunk &&
        (chunk.found ? (
          <>
            {chunk.chunks !== undefined && <p className={styles.meta}>Чанков на странице: {chunk.chunks}</p>}
            <pre className={styles.text}>{chunk.text || "— пусто —"}</pre>
          </>
        ) : (
          <Alert tone="warning">Страница {chunk.requested} не найдена в индексе</Alert>
        ))}
    </Widget>
  );
}
