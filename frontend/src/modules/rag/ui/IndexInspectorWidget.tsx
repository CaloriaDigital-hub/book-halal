"use client";

import { useState } from "react";
import { Alert, Badge, Button, EmptyState, Spinner, Widget } from "@/shared/ui";
import { useAsync } from "@/shared/lib/use-async";
import { inspectIndex } from "../api/rag-api";
import styles from "./IndexInspectorWidget.module.css";

const LIMITS = [10, 50, 200];

/** `version` меняется снаружи, когда индекс пересобрали, — виджет перечитает данные. */
export function IndexInspectorWidget({ bookId, version = 0 }: { bookId: string; version?: number }) {
  const [limit, setLimit] = useState(LIMITS[0]);
  const { data, error, loading, reload } = useAsync(`index:${bookId}:${limit}:${version}`, () =>
    inspectIndex(bookId, limit),
  );

  return (
    <Widget
      title="Индекс в Qdrant"
      description="Чанки, по которым RAG ищет ответы"
      actions={
        <Button size="sm" onClick={reload} loading={loading && !!data}>
          Обновить
        </Button>
      }
    >
      {error && <Alert tone="danger" title="Не удалось прочитать индекс">{error.message}</Alert>}
      {!data && loading && <Spinner label="Загрузка…" />}
      {data && (
        <>
          <div className={styles.toolbar}>
            <Badge tone={data.total > 0 ? "success" : "warning"}>
              {data.total > 0 ? `${data.total} чанков` : "Индекса нет"}
            </Badge>
            {data.total > 0 && (
              <label className={styles.limit}>
                Показать
                <select value={limit} onChange={(e) => setLimit(Number(e.target.value))}>
                  {LIMITS.map((n) => (
                    <option key={n} value={n}>
                      {n}
                    </option>
                  ))}
                </select>
              </label>
            )}
          </div>
          {data.total === 0 ? (
            <EmptyState>Книга не проиндексирована — RAG не сможет по ней отвечать.</EmptyState>
          ) : (
            <ol className={styles.chunks}>
              {data.chunks.map((chunk, i) => (
                <li key={i} className={styles.chunk}>
                  <span className={styles.page}>
                    стр. {chunk.page ?? "?"}
                    {chunk.chunk != null && ` · чанк ${chunk.chunk + 1}`}
                  </span>
                  <p className={styles.text}>{chunk.text || <em>пустой текст</em>}</p>
                </li>
              ))}
            </ol>
          )}
        </>
      )}
    </Widget>
  );
}
