"use client";

import Link from "next/link";
import { useState } from "react";
import { Alert, Button, EmptyState, Input, Spinner } from "@/shared/ui";
import { useAsync } from "@/shared/lib/use-async";
import { getCatalog } from "../api/storefront-api";
import { Cover } from "./Cover";
import { PriceTag } from "./PriceTag";
import styles from "./CatalogView.module.css";

export function CatalogView({ bookHref }: { bookHref: (id: string) => string }) {
  const { data, error, loading, reload } = useAsync("store-catalog", getCatalog);
  const [query, setQuery] = useState("");

  const q = query.trim().toLowerCase();
  const books = (data ?? []).filter((b) => !q || b.title.toLowerCase().includes(q) || b.author.toLowerCase().includes(q));

  return (
    <div className={styles.page}>
      <section className={styles.hero}>
        <h1 className={styles.heroTitle}>Исламские книги онлайн</h1>
        <p className={styles.heroText}>
          Читайте прямо в браузере и задавайте вопросы по тексту книги — ИИ ответит со ссылками на страницы.
        </p>
        <Input
          className={styles.search}
          type="search"
          placeholder="Поиск по названию или автору"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          aria-label="Поиск книг"
        />
      </section>

      {error && (
        <Alert tone="danger" title="Не удалось загрузить каталог" action={<Button size="sm" onClick={reload}>Повторить</Button>}>
          {error.message}
        </Alert>
      )}
      {!data && loading && <Spinner label="Загружаем книги…" />}
      {data && books.length === 0 && <EmptyState>{q ? "Ничего не нашлось" : "Книг пока нет"}</EmptyState>}

      {books.length > 0 && (
        <ul className={styles.grid}>
          {books.map((book) => (
            <li key={book.id}>
              <Link href={bookHref(book.id)} className={styles.card}>
                <Cover src={book.cover_url} title={book.title} />
                <div className={styles.cardBody}>
                  <span className={styles.title}>{book.title}</span>
                  <span className={styles.author}>{book.author}</span>
                  <PriceTag price={book.price} />
                </div>
              </Link>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
