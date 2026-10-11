"use client";

import { useParams } from "next/navigation";
import { ReaderView } from "@/modules/reader";
import { useStoreBook } from "@/modules/storefront";
import { storeRoutes } from "../../../../routes";

export function ReadScreen() {
  const { id } = useParams<{ id: string }>();
  const { data: book } = useStoreBook(id);

  return (
    <ReaderView bookId={id} title={book?.title ?? ""} backHref={storeRoutes.book(id)} buyHref={storeRoutes.buy(id)} />
  );
}
