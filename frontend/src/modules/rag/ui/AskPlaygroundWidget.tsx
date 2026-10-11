"use client";

import { useState, type FormEvent } from "react";
import { Alert, Button, Field, Textarea, Widget } from "@/shared/ui";
import { askBook } from "../api/rag-api";
import type { AskResult } from "../model/types";
import styles from "./AskPlaygroundWidget.module.css";

export function AskPlaygroundWidget({ bookId }: { bookId: string }) {
  const [question, setQuestion] = useState("");
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [result, setResult] = useState<(AskResult & { tookMs: number }) | null>(null);

  async function onSubmit(event: FormEvent) {
    event.preventDefault();
    if (!question.trim()) return;
    setPending(true);
    setError(null);
    const started = performance.now();
    try {
      const res = await askBook(bookId, question);
      setResult({ ...res, tookMs: Math.round(performance.now() - started) });
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    } finally {
      setPending(false);
    }
  }

  return (
    <Widget title="Проверить ответы" description="Задайте вопрос по книге — так же, как это делает читатель">
      <form className={styles.form} onSubmit={onSubmit}>
        <Field label="Вопрос">
          <Textarea
            value={question}
            onChange={(e) => setQuestion(e.target.value)}
            placeholder="О чём эта книга?"
            onKeyDown={(e) => {
              if (e.key === "Enter" && (e.ctrlKey || e.metaKey)) e.currentTarget.form?.requestSubmit();
            }}
          />
        </Field>
        <div className={styles.submit}>
          <Button type="submit" variant="primary" loading={pending} disabled={!question.trim()}>
            Спросить
          </Button>
          <span className={styles.hint}>Ctrl+Enter</span>
        </div>
      </form>

      {error && <Alert tone="danger" title="RAG не ответил">{error}</Alert>}
      {result && (
        <div className={styles.result}>
          <div className={styles.answerHeader}>
            <h3>Ответ</h3>
            <span className={styles.hint}>{(result.tookMs / 1000).toFixed(1)} с</span>
          </div>
          <p className={styles.answer}>{result.answer}</p>
          {result.sources.length > 0 && (
            <>
              <h3>Источники</h3>
              <ul className={styles.sources}>
                {result.sources.map((source, i) => (
                  <li key={i}>
                    <span className={styles.page}>стр. {source.page}</span>
                    <p className={styles.sourceText}>{source.text}</p>
                  </li>
                ))}
              </ul>
            </>
          )}
        </div>
      )}
    </Widget>
  );
}
