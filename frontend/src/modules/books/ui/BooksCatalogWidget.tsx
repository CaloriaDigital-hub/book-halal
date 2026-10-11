"use client";

import Link from "next/link";
import { useState } from "react";
import { Alert, Button, EmptyState, Input, Spinner, Widget } from "@/shared/ui";
import { useAsync } from "@/shared/lib/use-async";
import { getCatalog } from "../api/books-api";
import { formatPrice } from "../model/status";
import { BookCover } from "./BookCover";
import styles from "./BooksCatalogWidget.module.css";

export function BooksCatalogWidget({ bookHref }: { bookHref: (id: string) => string }) {
  const { data, error, loading, reload } = useAsync("catalog", getCatalog);
  const [query, setQuery] = useState("");

  const q = query.trim().toLowerCase();
  const books = (data ?? []).filter(
    (b) => !q || b.title.toLowerCase().includes(q) || b.author.toLowerCase().includes(q) || b.id.includes(q),
  );

  return (
    <Widget
      title="Каталог"
      description="Только книги в статусе «готова». Загруженные недавно появятся после обработки."
      actions={
        <Button size="sm" onClick={reload} loading={loading && !!data}>
          Обновить
        </Button>
      }
    >
      <div className={styles.toolbar}>
        <Input placeholder="Поиск по названию, автору или ID" value={query} onChange={(e) => setQuery(e.target.value)} />
        {data && <span className={styles.count}>{books.length} из {data.length}</span>}
      </div>

      {error && <Alert tone="danger" title="Не удалось загрузить каталог">{error.message}</Alert>}
      {!data && loading && <Spinner label="Загрузка…" />}
      {data && books.length === 0 && <EmptyState>{q ? "Ничего не найдено" : "Каталог пуст"}</EmptyState>}

      {books.length > 0 && (
        <ul className={styles.list}>
          {books.map((book) => (
            <li key={book.id}>
              <Link href={bookHref(book.id)} className={styles.row}>
                <BookCover src={book.cover_url} title={book.title} />
                <div className={styles.info}>
                  <span className={styles.title}>{book.title}</span>
                  <span className={styles.meta}>{book.author}</span>
                </div>
                <span className={styles.price}>{formatPrice(book.price)}</span>
              </Link>
            </li>
          ))}
        </ul>
      )}
    </Widget>
  );
}
