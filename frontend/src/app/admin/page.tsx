"use client";

import { CatalogStatsWidget } from "@/modules/books";
import { RagHealthWidget } from "@/modules/rag";
import { Grid, PageHeader } from "@/shared/ui";
import { adminRoutes } from "./routes";

export default function OverviewPage() {
  return (
    <>
      <PageHeader title="Обзор" subtitle="Состояние сервисов и каталога" />
      <Grid columns={2}>
        <RagHealthWidget />
        <CatalogStatsWidget catalogHref={adminRoutes.books} />
      </Grid>
    </>
  );
}
