"use client";

import { useState, type FormEvent } from "react";
import { Alert, Button, Spinner, Textarea } from "@/shared/ui";
import { ask } from "../api/reader-api";
import type { AskAnswer } from "../model/types";
import styles from "./AskPanel.module.css";

type Turn = { question: string; answer?: AskAnswer; error?: string };

const EXAMPLES = ["О чём эта книга?", "Как совершать омовение?", "Что нарушает намаз?"];

/** Вопросы по тексту книги. Источники в ответе — ссылки на страницы. */
export function AskPanel({ bookId, onOpenPage }: { bookId: string; onOpenPage: (page: number) => void }) {
  const [turns, setTurns] = useState<Turn[]>([]);
  const [question, setQuestion] = useState("");
  const [pending, setPending] = useState(false);

  async function send(q: string) {
    const text = q.trim();
    if (!text || pending) return;
    setQuestion("");
    setPending(true);
    setTurns((t) => [...t, { question: text }]);
    try {
      const answer = await ask(bookId, text);
      setTurns((t) => t.map((turn, i) => (i === t.length - 1 ? { ...turn, answer } : turn)));
    } catch (e) {
      const error = e instanceof Error ? e.message : String(e);
      setTurns((t) => t.map((turn, i) => (i === t.length - 1 ? { ...turn, error } : turn)));
    } finally {
      setPending(false);
    }
  }

  function onSubmit(event: FormEvent) {
    event.preventDefault();
    send(question);
  }

  return (
    <div className={styles.panel}>
      <div className={styles.thread}>
        {turns.length === 0 && (
          <div className={styles.intro}>
            <p>Спросите что-нибудь по книге — ответ будет со ссылками на страницы.</p>
            <div className={styles.examples}>
              {EXAMPLES.map((e) => (
                <button key={e} type="button" className={styles.example} onClick={() => send(e)}>
                  {e}
                </button>
              ))}
            </div>
          </div>
        )}
        {turns.map((turn, i) => (
          <div key={i} className={styles.turn}>
            <p className={styles.question}>{turn.question}</p>
            {!turn.answer && !turn.error && <Spinner label="Ищу ответ в книге…" />}
            {turn.error && <Alert tone="danger">{turn.error}</Alert>}
            {turn.answer && (
              <div className={styles.answer}>
                <p className={styles.answerText}>{turn.answer.answer}</p>
                {!!turn.answer.sources?.length && (
                  <div className={styles.sources}>
                    {uniquePages(turn.answer.sources.map((s) => s.page)).map((page) => (
                      <button key={page} type="button" className={styles.source} onClick={() => onOpenPage(page)}>
                        стр. {page}
                      </button>
                    ))}
                  </div>
                )}
              </div>
            )}
          </div>
        ))}
      </div>

      <form className={styles.form} onSubmit={onSubmit}>
        <Textarea
          value={question}
          onChange={(e) => setQuestion(e.target.value)}
          placeholder="Ваш вопрос…"
          rows={2}
          onKeyDown={(e) => {
            if (e.key === "Enter" && !e.shiftKey) {
              e.preventDefault();
              send(question);
            }
          }}
        />
        <Button type="submit" variant="primary" loading={pending} disabled={!question.trim()}>
          Спросить
        </Button>
      </form>
    </div>
  );
}

function uniquePages(pages: number[]) {
  return [...new Set(pages)].sort((a, b) => a - b);
}
