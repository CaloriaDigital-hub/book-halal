"use client";

import { useState, type FormEvent } from "react";
import { Alert, Button, Field, Input, Spinner, Textarea, Widget } from "@/shared/ui";
import { useAsync } from "@/shared/lib/use-async";
import { getRequisites, updateRequisites } from "../api/payments-api";
import type { PaymentSettings } from "../model/types";
import styles from "./RequisitesWidget.module.css";

export function RequisitesWidget() {
  const { data, error, loading } = useAsync("requisites", getRequisites);

  return (
    <Widget title="Реквизиты для оплаты" description="Покупатель видит их перед тем, как отправить чек">
      {error && <Alert tone="danger" title="Не удалось загрузить реквизиты">{error.message}</Alert>}
      {!data && loading && <Spinner label="Загрузка…" />}
      {/* key: после загрузки форма стартует с актуальными данными */}
      {data && <RequisitesForm key={data.updated_at} initial={data} />}
    </Widget>
  );
}

function RequisitesForm({ initial }: { initial: PaymentSettings }) {
  const [settings, setSettings] = useState(initial);
  const [details, setDetails] = useState(initial.details);
  const [qr, setQr] = useState<File | null>(null);
  const [fileKey, setFileKey] = useState(0);
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [saved, setSaved] = useState(false);

  async function save(removeQr = false) {
    setPending(true);
    setError(null);
    setSaved(false);
    try {
      const next = await updateRequisites({ details, qr: removeQr ? null : qr, removeQr });
      setSettings(next);
      setQr(null);
      setFileKey((k) => k + 1);
      setSaved(true);
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    } finally {
      setPending(false);
    }
  }

  function onSubmit(event: FormEvent) {
    event.preventDefault();
    save();
  }

  return (
    <form className={styles.form} onSubmit={onSubmit}>
      <Field label="Реквизиты" hint="Банк, номер карты или счёта, получатель — как увидит покупатель">
        <Textarea
          value={details}
          onChange={(e) => setDetails(e.target.value)}
          placeholder={"MBank: 0555 123 456\nПолучатель: Иван И."}
          rows={4}
        />
      </Field>

      <div className={styles.qrRow}>
        <div className={styles.qrPreview}>
          {settings.qr_url ? (
            // eslint-disable-next-line @next/next/no-img-element -- QR отдаёт Go-бэк из /static
            <img src={settings.qr_url} alt="QR для оплаты" />
          ) : (
            <span>QR нет</span>
          )}
        </div>
        <div className={styles.qrControls}>
          <Field label={settings.qr_url ? "Заменить QR-код" : "QR-код"} hint="jpg, png или webp до 5 МБ">
            <Input key={fileKey} type="file" accept="image/jpeg,image/png,image/webp" onChange={(e) => setQr(e.target.files?.[0] ?? null)} />
          </Field>
          {settings.qr_url && (
            <div>
              <Button variant="ghost" size="sm" disabled={pending} onClick={() => save(true)}>
                Убрать QR
              </Button>
            </div>
          )}
        </div>
      </div>

      {error && <Alert tone="danger">{error}</Alert>}
      {saved && <Alert tone="success">Сохранено</Alert>}

      <div>
        <Button type="submit" variant="primary" loading={pending}>
          Сохранить
        </Button>
      </div>
    </form>
  );
}
