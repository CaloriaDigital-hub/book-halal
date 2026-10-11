"use client";

import Link from "next/link";
import { useState, type FormEvent } from "react";
import { formatSom } from "@/shared/lib/format";
import { useAsync } from "@/shared/lib/use-async";
import { Alert, Button, CopyButton, Field, Input, Segmented, Spinner } from "@/shared/ui";
import { getRequisites, submitPurchase } from "../api/purchase-api";
import { humanize } from "../model/messages";
import type { ContactType } from "../model/types";
import { useMyPurchase } from "../model/use-my-purchase";
import { ReceiptPicker } from "./ReceiptPicker";
import styles from "./PurchaseFlow.module.css";

type Props = {
  bookId: string;
  bookTitle: string;
  price: number;
  bookHref: string;
  readHref: string;
};

const CONTACTS: { value: ContactType; label: string; placeholder: string; inputMode: "tel" | "text" }[] = [
  { value: "whatsapp", label: "WhatsApp", placeholder: "+996 555 123 456", inputMode: "tel" },
  { value: "telegram", label: "Telegram", placeholder: "@username или номер", inputMode: "text" },
  { value: "phone", label: "Телефон", placeholder: "+996 555 123 456", inputMode: "tel" },
];

/** Оплата: 1) перевод по QR/реквизитам → 2) чек и контакт → 3) статус проверки. */
export function PurchaseFlow({ bookId, bookTitle, price, bookHref, readHref }: Props) {
  const { data: purchase, loading, reload } = useMyPurchase(bookId, true);
  const [step, setStep] = useState<"pay" | "receipt">("pay");

  if (loading && purchase === undefined) return <Spinner label="Загрузка…" />;

  if (purchase?.status === "pending" || purchase?.status === "approved") {
    return <StatusCard status={purchase.status} bookHref={bookHref} readHref={readHref} />;
  }

  return (
    <div className={styles.flow}>
      <header className={styles.header}>
        <p className={styles.kicker}>Покупка книги</p>
        <h1 className={styles.title}>{bookTitle}</h1>
        <p className={styles.amount}>{formatSom(price)}</p>
      </header>

      <ol className={styles.steps}>
        <li className={step === "pay" ? styles.current : styles.done}>Оплата</li>
        <li className={step === "receipt" ? styles.current : undefined}>Чек</li>
        <li>Проверка</li>
      </ol>

      {purchase?.status === "rejected" && step === "pay" && (
        <Alert tone="warning" title="Прошлый чек не подтвердили">
          Проверьте, что сумма и реквизиты совпадают, и отправьте чек ещё раз.
        </Alert>
      )}

      {step === "pay" ? (
        <PayStep price={price} onPaid={() => setStep("receipt")} />
      ) : (
        <ReceiptStep bookId={bookId} onBack={() => setStep("pay")} onSent={reload} />
      )}
    </div>
  );
}

function PayStep({ price, onPaid }: { price: number; onPaid: () => void }) {
  const { data, error, loading } = useAsync("requisites", getRequisites);

  if (error) return <Alert tone="danger" title="Не удалось загрузить реквизиты">{error.message}</Alert>;
  if (!data || loading) return <Spinner label="Загружаем реквизиты…" />;

  const configured = data.details.trim() !== "" || data.qr_url !== "";
  if (!configured) {
    return (
      <Alert tone="warning" title="Оплата временно недоступна">
        Реквизиты ещё не настроены. Попробуйте позже.
      </Alert>
    );
  }

  return (
    <section className={styles.card}>
      <h2 className={styles.cardTitle}>Переведите {formatSom(price)}</h2>
      <div className={styles.payLayout}>
        {data.qr_url && (
          <figure className={styles.qr}>
            {/* eslint-disable-next-line @next/next/no-img-element -- QR из /static Go-бэка */}
            <img src={data.qr_url} alt="QR-код для оплаты" />
            <figcaption>Отсканируйте в приложении банка</figcaption>
          </figure>
        )}
        {data.details && (
          <div className={styles.requisites}>
            <p className={styles.label}>{data.qr_url ? "Или по реквизитам" : "Реквизиты"}</p>
            <pre className={styles.details}>{data.details}</pre>
            <CopyButton text={data.details} label="Скопировать реквизиты" />
          </div>
        )}
      </div>
      <Alert tone="info">
        Переводите ровно <b>{formatSom(price)}</b>. После перевода сохраните скриншот или PDF чека — он понадобится на следующем шаге.
      </Alert>
      <Button variant="primary" onClick={onPaid} className={styles.next}>
        Я оплатил(а) — прикрепить чек
      </Button>
    </section>
  );
}

function ReceiptStep({ bookId, onBack, onSent }: { bookId: string; onBack: () => void; onSent: () => void }) {
  const [receipt, setReceipt] = useState<File | null>(null);
  const [contactType, setContactType] = useState<ContactType>("whatsapp");
  const [contact, setContact] = useState("");
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const current = CONTACTS.find((c) => c.value === contactType)!;

  async function onSubmit(event: FormEvent) {
    event.preventDefault();
    if (!receipt) return setError("Прикрепите чек");
    setPending(true);
    setError(null);
    try {
      await submitPurchase(bookId, { receipt, contactType, contact });
      onSent();
    } catch (e) {
      setError(humanize(e));
      setPending(false);
    }
  }

  return (
    <form className={styles.card} onSubmit={onSubmit}>
      <h2 className={styles.cardTitle}>Прикрепите чек</h2>
      <ReceiptPicker
        file={receipt}
        onChange={(f) => {
          setReceipt(f);
          setError(null);
        }}
        onError={setError}
      />

      <div className={styles.contact}>
        <p className={styles.label}>Как с вами связаться, если с оплатой что-то не так?</p>
        <Segmented label="Способ связи" options={CONTACTS} value={contactType} onChange={setContactType} />
        <Field label={current.label}>
          <Input
            required
            inputMode={current.inputMode}
            placeholder={current.placeholder}
            value={contact}
            onChange={(e) => setContact(e.target.value)}
          />
        </Field>
      </div>

      {error && <Alert tone="danger">{error}</Alert>}

      <div className={styles.buttons}>
        <Button variant="ghost" onClick={onBack} disabled={pending}>
          ← Назад
        </Button>
        <Button type="submit" variant="primary" loading={pending}>
          Отправить чек
        </Button>
      </div>
    </form>
  );
}

function StatusCard({ status, bookHref, readHref }: { status: "pending" | "approved"; bookHref: string; readHref: string }) {
  if (status === "approved") {
    return (
      <section className={`${styles.card} ${styles.status}`}>
        <span className={styles.statusIcon} aria-hidden>
          ✅
        </span>
        <h1 className={styles.cardTitle}>Оплата подтверждена</h1>
        <p className={styles.statusText}>Книга ваша — читайте и задавайте вопросы по тексту.</p>
        <Link href={readHref} className={styles.readLink}>
          Начать читать
        </Link>
      </section>
    );
  }
  return (
    <section className={`${styles.card} ${styles.status}`}>
      <span className={styles.statusIcon} aria-hidden>
        ⏳
      </span>
      <h1 className={styles.cardTitle}>Чек отправлен</h1>
      <p className={styles.statusText}>
        Проверяем оплату. Как только подтвердим, книга откроется — эта страница обновится сама.
      </p>
      <Spinner label="Ждём подтверждения…" />
      <Link href={bookHref} className={styles.backLink}>
        Вернуться к книге
      </Link>
    </section>
  );
}
