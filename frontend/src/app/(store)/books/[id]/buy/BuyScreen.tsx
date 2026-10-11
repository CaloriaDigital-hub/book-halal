"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
import { PurchaseFlow } from "@/modules/purchase";
import { useStoreBook } from "@/modules/storefront";
import { Alert, Spinner } from "@/shared/ui";
import { storeRoutes } from "../../../../routes";

export function BuyScreen() {
  const { id } = useParams<{ id: string }>();
  const { data: book, error } = useStoreBook(id);

  if (error) return <Alert tone="danger" title="Книга не найдена">{error.message}</Alert>;
  if (!book) return <Spinner label="Загрузка…" />;

  if (book.has_access || book.price === 0) {
    return (
      <Alert tone="success" title="Книга уже доступна">
        <Link href={storeRoutes.read(book.id)}>Читать →</Link>
      </Alert>
    );
  }

  return (
    <PurchaseFlow
      bookId={book.id}
      bookTitle={book.title}
      price={book.price}
      bookHref={storeRoutes.book(book.id)}
      readHref={storeRoutes.read(book.id)}
    />
  );
}
