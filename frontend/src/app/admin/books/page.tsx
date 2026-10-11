"use client";

import { BooksCatalogWidget, UploadBookWidget } from "@/modules/books";
import { Grid, PageHeader } from "@/shared/ui";
import { adminRoutes } from "../routes";

export default function BooksPage() {
  return (
    <>
      <PageHeader title="Книги" subtitle="Загрузка PDF и каталог" />
      <Grid columns={2}>
        <UploadBookWidget bookHref={adminRoutes.book} />
        <BooksCatalogWidget bookHref={adminRoutes.book} />
      </Grid>
    </>
  );
}
