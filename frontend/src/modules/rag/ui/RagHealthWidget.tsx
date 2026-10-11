"use client";

import { useEffect } from "react";
import { Badge, Button, Spinner, Widget } from "@/shared/ui";
import { useAsync } from "@/shared/lib/use-async";
import { getHealth } from "../api/rag-api";
import styles from "./RagHealthWidget.module.css";

const POLL_MS = 15_000;

async function probe() {
  const started = performance.now();
  const health = await getHealth();
  return { ...health, latencyMs: Math.round(performance.now() - started), checkedAt: new Date() };
}

export function RagHealthWidget() {
  const { data, error, loading, reload } = useAsync("rag-health", probe);

  useEffect(() => {
    const timer = setInterval(reload, POLL_MS);
    return () => clearInterval(timer);
  }, [reload]);

  const ok = !error && data?.status === "ok";

  return (
    <Widget
      title="RAG-сервис"
      description="Эмбеддинги, реранкер и генерация ответов"
      actions={
        <Button size="sm" onClick={reload} loading={loading && !!(data || error)}>
          Проверить
        </Button>
      }
    >
      {!data && !error ? (
        <Spinner label="Проверяю…" />
      ) : (
        <div className={styles.status}>
          <Badge tone={ok ? "success" : "danger"}>{ok ? "Работает" : "Недоступен"}</Badge>
          {ok && data && <span className={styles.meta}>ответ за {data.latencyMs} мс</span>}
          {error && <span className={styles.error}>{error.message}</span>}
        </div>
      )}
      {data && <p className={styles.checked}>Проверено в {data.checkedAt.toLocaleTimeString("ru-RU")}, автообновление раз в 15 с</p>}
    </Widget>
  );
}
