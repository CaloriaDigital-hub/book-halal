"use client";

import { useEffect, useState } from "react";
import { Alert, Button, Widget } from "@/shared/ui";
import { deleteIndex } from "../api/rag-api";
import type { IndexPageInput } from "../model/types";
import { useReindex, type ReindexState } from "../model/use-reindex";
import styles from "./ReindexWidget.module.css";

type Props = {
  bookId: string;
  /** Откуда брать страницы — решает тот, кто собирает экран (rag не знает про books). */
  loadPages: () => Promise<IndexPageInput[]>;
  onChanged?: () => void;
};

export function ReindexWidget({ bookId, loadPages, onChanged }: Props) {
  const { state, running, start } = useReindex(bookId, loadPages, onChanged);
  const [deleting, setDeleting] = useState(false);
  const [deleteError, setDeleteError] = useState<string | null>(null);

  async function onReindex() {
    if (!confirm("Текущий индекс будет удалён и собран заново через OCR всех страниц. Продолжить?")) return;
    setDeleteError(null);
    await start();
  }

  async function onDelete() {
    if (!confirm("Удалить индекс книги? RAG перестанет отвечать по ней до переиндексации.")) return;
    setDeleting(true);
    setDeleteError(null);
    try {
      await deleteIndex(bookId);
      onChanged?.();
    } catch (e) {
      setDeleteError(e instanceof Error ? e.message : String(e));
    } finally {
      setDeleting(false);
    }
  }

  return (
    <Widget title="Переиндексация" description="OCR всех страниц → чанки → Qdrant. Не закрывайте вкладку до конца.">
      <div className={styles.actions}>
        <Button variant="primary" onClick={onReindex} loading={running} disabled={deleting}>
          {running ? "Индексирую…" : "Переиндексировать"}
        </Button>
        <Button variant="danger" onClick={onDelete} loading={deleting} disabled={running}>
          Удалить индекс
        </Button>
      </div>

      <Progress state={state} />
      {deleteError && <Alert tone="danger" title="Не удалось удалить индекс">{deleteError}</Alert>}
    </Widget>
  );
}

const PHASE_LABEL: Partial<Record<ReindexState["phase"], string>> = {
  "loading-pages": "Получаю список страниц…",
  deleting: "Удаляю старый индекс…",
  indexing: "Распознаю и индексирую страницы — это может занять несколько минут…",
};

function Progress({ state }: { state: ReindexState }) {
  if (state.phase === "idle") return null;

  const chunks = state.chunks !== undefined ? ` (${state.chunks} чанков)` : "";
  const counters = `${state.total} стр. · в индексе ${state.indexed}${chunks} · пропущено ${state.skipped}`;

  return (
    <div className={styles.progress}>
      {PHASE_LABEL[state.phase] && (
        <>
          <p className={styles.phase}>{PHASE_LABEL[state.phase]}</p>
          <div className={styles.bar} role="progressbar" aria-busy>
            <div className={styles.indeterminate} />
          </div>
          {state.phase === "indexing" && state.startedAt && (
            <p className={styles.counters}>
              {state.total} стр. · идёт <Elapsed since={state.startedAt} />
            </p>
          )}
        </>
      )}
      {state.phase === "done" && (
        <Alert tone={state.skipped ? "warning" : "success"} title="Готово">
          {counters}
          {state.skipped > 0 && ". Пропущенные страницы — пустой результат OCR."}
        </Alert>
      )}
      {state.phase === "error" && (
        <Alert tone="danger" title="Переиндексация упала">
          {state.error}
          {state.startedAt && ". Старый индекс уже удалён — книга без индекса, пока переиндексация не пройдёт."}
        </Alert>
      )}
    </div>
  );
}

function Elapsed({ since }: { since: number }) {
  const [now, setNow] = useState(since);
  useEffect(() => {
    const timer = setInterval(() => setNow(Date.now()), 1000);
    return () => clearInterval(timer);
  }, []);
  const sec = Math.max(0, Math.round((now - since) / 1000));
  return (
    <>
      {Math.floor(sec / 60)}:{String(sec % 60).padStart(2, "0")}
    </>
  );
}
