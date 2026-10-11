"use client";

import { PurchasesWidget, RequisitesWidget } from "@/modules/payments";
import { Grid, PageHeader } from "@/shared/ui";

export default function PaymentsPage() {
  return (
    <>
      <PageHeader title="Оплаты" subtitle="Чеки покупателей и реквизиты для оплаты" />
      <Grid columns={1}>
        <PurchasesWidget />
      </Grid>
      <Grid columns={2}>
        <RequisitesWidget />
      </Grid>
    </>
  );
}
