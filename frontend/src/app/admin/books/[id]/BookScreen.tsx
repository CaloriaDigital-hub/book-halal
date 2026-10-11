"use client";

import { useParams } from "next/navigation";
import { useCallback, useState } from "react";
import { BookSummaryWidget, getAllPages } from "@/modules/books";
import {
  AskPlaygroundWidget,
  IndexInspectorWidget,
  PageChunkWidget,
  ReindexWidget,
  type IndexPageInput,
} from "@/modules/rag";
import { toBackendPublicUrl } from "@/shared/config/public-env";
import { Grid, PageHeader } from "@/shared/ui";
import { adminRoutes } from "../../routes";

export function BookScreen() {
  const { id: bookId } = useParams<{ id: string }>();
  // Пересобрали/удалили индекс → инспектор перечитает данные.
  const [indexVersion, setIndexVersion] = useState(0);

  // Склейка модулей: rag получает страницы из books, не завися от него напрямую.
  const loadPages = useCallback(
    async (): Promise<IndexPageInput[]> =>
      (await getAllPages(bookId)).map((p) => ({ page_number: p.page_number, image_url: toBackendPublicUrl(p.image_url) })),
    [bookId],
  );
  const onIndexChanged = useCallback(() => setIndexVersion((v) => v + 1), []);

  return (
    <>
      <PageHeader title="Книга и RAG" back={{ href: adminRoutes.books, label: "Все книги" }} />
      <Grid columns={2}>
        <BookSummaryWidget bookId={bookId} />
        <ReindexWidget bookId={bookId} loadPages={loadPages} onChanged={onIndexChanged} />
      </Grid>
      <Grid columns={2}>
        <AskPlaygroundWidget bookId={bookId} />
        <PageChunkWidget bookId={bookId} />
      </Grid>
      <Grid columns={1}>
        <IndexInspectorWidget bookId={bookId} version={indexVersion} />
      </Grid>
    </>
  );
}
