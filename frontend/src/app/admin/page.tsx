"use client";

import { CatalogStatsWidget } from "@/modules/books";
import { PendingPaymentsWidget } from "@/modules/payments";
import { RagHealthWidget } from "@/modules/rag";
import { Grid, PageHeader } from "@/shared/ui";
import { adminRoutes } from "./routes";

export default function OverviewPage() {
  return (
    <>
      <PageHeader title="Обзор" subtitle="Состояние сервисов, каталога и оплат" />
      <Grid columns={3}>
        <RagHealthWidget />
        <CatalogStatsWidget catalogHref={adminRoutes.books} />
        <PendingPaymentsWidget paymentsHref={adminRoutes.payments} />
      </Grid>
    </>
  );
}
