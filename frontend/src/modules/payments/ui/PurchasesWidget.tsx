"use client";

import { useEffect, useState } from "react";
import { formatDateTime, formatSom } from "@/shared/lib/format";
import { useAsync } from "@/shared/lib/use-async";
import { Alert, Badge, Button, EmptyState, Spinner, Widget, type Tone } from "@/shared/ui";
import { listPurchases, receiptUrl, reviewPurchase } from "../api/payments-api";
import { contactLink } from "../model/contact";
import type { AdminPurchase, PurchaseStatus } from "../model/types";
import styles from "./PurchasesWidget.module.css";

type Filter = PurchaseStatus | "all";

const FILTERS: { value: Filter; label: string }[] = [
  { value: "pending", label: "Ожидают" },
  { value: "approved", label: "Подтверждённые" },
  { value: "rejected", label: "Отклонённые" },
  { value: "all", label: "Все" },
];

const STATUS: Record<PurchaseStatus, { label: string; tone: Tone }> = {
  pending: { label: "Ожидает", tone: "warning" },
  approved: { label: "Подтверждена", tone: "success" },
  rejected: { label: "Отклонена", tone: "danger" },
};

const POLL_MS = 30_000;

export function PurchasesWidget() {
  const [filter, setFilter] = useState<Filter>("pending");
  const { data, error, loading, reload } = useAsync(`purchases:${filter}`, () => listPurchases(filter));

  useEffect(() => {
    const timer = setInterval(reload, POLL_MS);
    return () => clearInterval(timer);
  }, [reload]);

  return (
    <Widget
      title="Оплаты"
      description="Основное подтверждение — кнопками в Telegram. Здесь — то же самое, если бот недоступен."
      actions={
        <Button size="sm" onClick={reload} loading={loading && !!data}>
          Обновить
        </Button>
      }
    >
      <div className={styles.filters} role="tablist">
        {FILTERS.map((f) => (
          <button
            key={f.value}
            role="tab"
            aria-selected={filter === f.value}
            className={filter === f.value ? `${styles.filter} ${styles.filterActive}` : styles.filter}
            onClick={() => setFilter(f.value)}
          >
            {f.label}
          </button>
        ))}
      </div>

      {error && <Alert tone="danger" title="Не удалось загрузить оплаты">{error.message}</Alert>}
      {!data && loading && <Spinner label="Загрузка…" />}
      {data && data.length === 0 && <EmptyState>{filter === "pending" ? "Новых чеков нет" : "Пусто"}</EmptyState>}

      {data && data.length > 0 && (
        <ul className={styles.list}>
          {data.map((p) => (
            <PurchaseRow key={p.id} purchase={p} onReviewed={reload} />
          ))}
        </ul>
      )}
    </Widget>
  );
}

function PurchaseRow({ purchase: p, onReviewed }: { purchase: AdminPurchase; onReviewed: () => void }) {
  const [pending, setPending] = useState<"approve" | "reject" | null>(null);
  const [error, setError] = useState<string | null>(null);
  const contact = contactLink(p.contact_type, p.contact_value);
  const isImage = p.receipt_mime.startsWith("image/");

  async function review(approve: boolean) {
    if (!approve && !confirm(`Отклонить оплату «${p.book_title}» от ${p.buyer_name}?`)) return;
    setPending(approve ? "approve" : "reject");
    setError(null);
    try {
      await reviewPurchase(p.id, approve);
      onReviewed();
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
      setPending(null);
    }
  }

  return (
    <li className={styles.row}>
      <a href={receiptUrl(p.id)} target="_blank" rel="noreferrer" className={styles.receipt} title="Открыть чек">
        {isImage ? (
          // eslint-disable-next-line @next/next/no-img-element -- чек приходит из BFF с cookie, оптимизатор Next не нужен
          <img src={receiptUrl(p.id)} alt="Чек" loading="lazy" />
        ) : (
          <span>PDF</span>
        )}
      </a>

      <div className={styles.info}>
        <div className={styles.head}>
          <span className={styles.title}>{p.book_title}</span>
          <Badge tone={STATUS[p.status].tone}>{STATUS[p.status].label}</Badge>
        </div>
        <span className={styles.amount}>{formatSom(p.amount)}</span>
        <span className={styles.meta}>
          {p.buyer_name} · {p.buyer_email}
        </span>
        <span className={styles.meta}>
          {contact.label}:{" "}
          {contact.href ? (
            <a href={contact.href} target="_blank" rel="noreferrer">
              {p.contact_value}
            </a>
          ) : (
            p.contact_value
          )}
        </span>
        <span className={styles.meta}>
          {formatDateTime(p.created_at)}
          {p.reviewed_at && ` · ${p.status === "approved" ? "подтвердил" : "отклонил"} ${p.reviewed_by} в ${formatDateTime(p.reviewed_at)}`}
        </span>
        {error && <Alert tone="danger">{error}</Alert>}
      </div>

      {p.status === "pending" && (
        <div className={styles.actions}>
          <Button variant="primary" size="sm" loading={pending === "approve"} disabled={!!pending} onClick={() => review(true)}>
            Подтвердить
          </Button>
          <Button variant="danger" size="sm" loading={pending === "reject"} disabled={!!pending} onClick={() => review(false)}>
            Отклонить
          </Button>
        </div>
      )}
    </li>
  );
}
