"use client";

import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { useState, type FormEvent } from "react";
import { Alert, Button, Field, Input } from "@/shared/ui";
import { confirmRegistration, requestCode, signIn } from "../api/session-api";
import { safeNext } from "../model/safe-next";
import { setSession } from "../model/use-session";
import styles from "./AuthCard.module.css";

/** Регистрация в два шага: email → код из письма + имя и пароль. После неё сразу входим. */
export function RegisterForm({ loginHref }: { loginHref: string }) {
  const router = useRouter();
  const next = safeNext(useSearchParams().get("next"));
  const [step, setStep] = useState<"email" | "confirm">("email");
  const [form, setForm] = useState({ email: "", code: "", first_name: "", last_name: "", password: "" });
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const set = (key: keyof typeof form) => (e: { target: { value: string } }) =>
    setForm((prev) => ({ ...prev, [key]: e.target.value }));

  async function run(action: () => Promise<void>) {
    setPending(true);
    setError(null);
    try {
      await action();
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    } finally {
      setPending(false);
    }
  }

  function onSendCode(event: FormEvent) {
    event.preventDefault();
    run(async () => {
      await requestCode(form.email);
      setStep("confirm");
    });
  }

  function onConfirm(event: FormEvent) {
    event.preventDefault();
    run(async () => {
      await confirmRegistration({ ...form, code: form.code.trim() });
      const { role } = await signIn(form.email, form.password);
      setSession({ authenticated: true, role });
      router.replace(next ?? "/");
    });
  }

  const loginLink = next ? `${loginHref}?next=${encodeURIComponent(next)}` : loginHref;

  if (step === "email") {
    return (
      <form className={styles.card} onSubmit={onSendCode}>
        <div>
          <h1 className={styles.title}>Регистрация</h1>
          <p className={styles.subtitle}>Пришлём код подтверждения на почту</p>
        </div>
        {error && <Alert tone="danger">{error}</Alert>}
        <Field label="Email">
          <Input type="email" autoComplete="email" required value={form.email} onChange={set("email")} />
        </Field>
        <Button type="submit" variant="primary" loading={pending}>
          Получить код
        </Button>
        <p className={styles.switch}>
          Уже есть аккаунт? <Link href={loginLink}>Войти</Link>
        </p>
      </form>
    );
  }

  return (
    <form className={styles.card} onSubmit={onConfirm}>
      <div>
        <h1 className={styles.title}>Подтверждение</h1>
        <p className={styles.subtitle}>
          Код отправлен на <b>{form.email}</b>
        </p>
      </div>
      {error && <Alert tone="danger">{error}</Alert>}
      <Field label="Код из письма">
        <Input inputMode="numeric" autoComplete="one-time-code" required value={form.code} onChange={set("code")} />
      </Field>
      <div className={styles.row}>
        <Field label="Имя">
          <Input autoComplete="given-name" required value={form.first_name} onChange={set("first_name")} />
        </Field>
        <Field label="Фамилия">
          <Input autoComplete="family-name" required value={form.last_name} onChange={set("last_name")} />
        </Field>
      </div>
      <Field label="Пароль">
        <Input type="password" autoComplete="new-password" required minLength={8} value={form.password} onChange={set("password")} />
      </Field>
      <Button type="submit" variant="primary" loading={pending}>
        Создать аккаунт
      </Button>
      <p className={styles.switch}>
        <button type="button" className={styles.linkButton} onClick={() => setStep("email")}>
          Изменить email или отправить код снова
        </button>
      </p>
    </form>
  );
}
