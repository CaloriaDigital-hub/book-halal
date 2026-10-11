"use client";

import Link from "next/link";
import type { ReactNode } from "react";
import { Alert, Button, Spinner } from "@/shared/ui";
import { useAsync } from "@/shared/lib/use-async";
import { getBook } from "../api/storefront-api";
import type { BookDetails } from "../model/types";
import { Cover } from "./Cover";
import { PriceTag } from "./PriceTag";
import styles from "./BookView.module.css";

type Props = {
  bookId: string;
  backHref: string;
  /** Кнопки действий (купить/читать) — их собирает экран из других модулей. */
  actions: (book: BookDetails) => ReactNode;
};

export function BookView({ bookId, backHref, actions }: Props) {
  const { data: book, error, loading, reload } = useAsync(`store-book:${bookId}`, () => getBook(bookId));

  if (error) {
    return (
      <Alert tone="danger" title="Книга не найдена или недоступна" action={<Button size="sm" onClick={reload}>Повторить</Button>}>
        {error.message}
      </Alert>
    );
  }
  if (!book || loading) return <Spinner label="Загружаем книгу…" />;

  return (
    <article className={styles.book}>
      <Link href={backHref} className={styles.back}>
        ← Все книги
      </Link>
      <div className={styles.layout}>
        <Cover src={book.cover_url} title={book.title} className={styles.cover} />
        <div className={styles.info}>
          <h1 className={styles.title}>{book.title}</h1>
          <p className={styles.author}>{book.author}</p>
          <div className={styles.meta}>
            <PriceTag price={book.price} large />
            {book.total_pages > 0 && <span className={styles.pages}>{book.total_pages} стр.</span>}
          </div>
          <div className={styles.actions}>{actions(book)}</div>
          {book.description && <p className={styles.description}>{book.description}</p>}
        </div>
      </div>
    </article>
  );
}
