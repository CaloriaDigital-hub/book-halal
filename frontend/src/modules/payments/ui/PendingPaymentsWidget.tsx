"use client";

import Link from "next/link";
import { Alert, Spinner, Widget } from "@/shared/ui";
import { useAsync } from "@/shared/lib/use-async";
import { listPurchases } from "../api/payments-api";
import styles from "./PendingPaymentsWidget.module.css";

export function PendingPaymentsWidget({ paymentsHref }: { paymentsHref: string }) {
  const { data, error } = useAsync("pending-payments", () => listPurchases("pending"));

  return (
    <Widget title="Оплаты" actions={<Link href={paymentsHref} className={styles.link}>Все оплаты →</Link>}>
      {error && <Alert tone="danger" title="Бэкенд не отвечает">{error.message}</Alert>}
      {!data && !error && <Spinner label="Загрузка…" />}
      {data && (
        <div>
          <p className={data.length > 0 ? `${styles.value} ${styles.attention}` : styles.value}>{data.length}</p>
          <p className={styles.label}>{data.length > 0 ? "чеков ждут проверки" : "новых чеков нет"}</p>
        </div>
      )}
    </Widget>
  );
}
