"use client";

import Link from "next/link";
import { useState, type FormEvent } from "react";
import { Alert, Button, Field, Input, Textarea, Widget } from "@/shared/ui";
import { uploadBook } from "../api/books-api";
import type { UploadBookResult } from "../model/types";
import styles from "./UploadBookWidget.module.css";

const EMPTY = { title: "", author: "", description: "", price: "0" };

export function UploadBookWidget({ bookHref }: { bookHref: (id: string) => string }) {
  const [form, setForm] = useState(EMPTY);
  const [file, setFile] = useState<File | null>(null);
  const [fileInputKey, setFileInputKey] = useState(0);
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [result, setResult] = useState<UploadBookResult | null>(null);

  const set = (key: keyof typeof EMPTY) => (e: { target: { value: string } }) =>
    setForm((prev) => ({ ...prev, [key]: e.target.value }));

  async function onSubmit(event: FormEvent) {
    event.preventDefault();
    if (!file) return;
    setPending(true);
    setError(null);
    setResult(null);
    try {
      const res = await uploadBook({ ...form, price: Number(form.price), document: file });
      setResult(res);
      setForm(EMPTY);
      setFile(null);
      setFileInputKey((k) => k + 1);
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    } finally {
      setPending(false);
    }
  }

  return (
    <Widget title="Загрузить книгу" description="PDF режется на страницы и автоматически отправляется в RAG на индексацию.">
      <form className={styles.form} onSubmit={onSubmit}>
        <div className={styles.row}>
          <Field label="Название">
            <Input required value={form.title} onChange={set("title")} />
          </Field>
          <Field label="Автор">
            <Input required value={form.author} onChange={set("author")} />
          </Field>
        </div>
        <Field label="Описание">
          <Textarea value={form.description} onChange={set("description")} />
        </Field>
        <div className={styles.row}>
          <Field label="Цена, сом" hint="0 — бесплатно">
            <Input type="number" min={0} step={1} required value={form.price} onChange={set("price")} />
          </Field>
          <Field label="PDF-файл" hint="До 100 МБ">
            <Input
              key={fileInputKey}
              type="file"
              accept="application/pdf,.pdf"
              required
              onChange={(e) => setFile(e.target.files?.[0] ?? null)}
            />
          </Field>
        </div>

        {error && <Alert tone="danger" title="Загрузка не удалась">{error}</Alert>}
        {result && (
          <Alert tone="success" title="Книга принята в обработку">
            <Link href={bookHref(result.book_id)}>Открыть книгу →</Link>
          </Alert>
        )}

        <div>
          <Button type="submit" variant="primary" loading={pending} disabled={!file}>
            Загрузить
          </Button>
        </div>
      </form>
    </Widget>
  );
}
