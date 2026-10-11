"use client";

import { Alert, Badge, Button, Spinner, Widget } from "@/shared/ui";
import { useAsync } from "@/shared/lib/use-async";
import { getBook } from "../api/books-api";
import { bookStatus, formatPrice } from "../model/status";
import { BookCover } from "./BookCover";
import styles from "./BookSummaryWidget.module.css";

export function BookSummaryWidget({ bookId }: { bookId: string }) {
  const { data: book, error, loading, reload } = useAsync(`book:${bookId}`, () => getBook(bookId));

  return (
    <Widget
      title="Книга"
      actions={
        <Button size="sm" onClick={reload} loading={loading && !!book}>
          Обновить
        </Button>
      }
    >
      {error && <Alert tone="danger" title="Не удалось загрузить книгу">{error.message}</Alert>}
      {!book && loading && <Spinner label="Загрузка…" />}
      {book && (
        <div className={styles.summary}>
          <BookCover src={book.cover_url} title={book.title} size="lg" />
          <dl className={styles.props}>
            <dt>Название</dt>
            <dd className={styles.strong}>{book.title}</dd>
            <dt>Автор</dt>
            <dd>{book.author}</dd>
            <dt>Статус</dt>
            <dd>
              <Badge tone={bookStatus(book.status).tone}>{bookStatus(book.status).label}</Badge>
            </dd>
            <dt>Страниц</dt>
            <dd>{book.total_pages}</dd>
            <dt>Цена</dt>
            <dd>{formatPrice(book.price)}</dd>
            <dt>ID</dt>
            <dd className={styles.mono}>{book.id}</dd>
            {book.description && (
              <>
                <dt>Описание</dt>
                <dd className={styles.description}>{book.description}</dd>
              </>
            )}
          </dl>
        </div>
      )}
    </Widget>
  );
}
