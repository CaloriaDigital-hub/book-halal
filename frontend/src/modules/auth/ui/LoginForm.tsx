"use client";

import { useRouter, useSearchParams } from "next/navigation";
import { useState, type FormEvent } from "react";
import { Alert, Button, Field, Input } from "@/shared/ui";
import { signIn } from "../api/session-api";
import styles from "./LoginForm.module.css";

export function LoginForm() {
  const router = useRouter();
  const next = safeNext(useSearchParams().get("next"));
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function onSubmit(event: FormEvent) {
    event.preventDefault();
    setPending(true);
    setError(null);
    try {
      await signIn(email, password);
      router.replace(next);
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
      setPending(false);
    }
  }

  return (
    <form className={styles.card} onSubmit={onSubmit}>
      <div>
        <h1 className={styles.title}>Book Halal</h1>
        <p className={styles.subtitle}>Вход в админку</p>
      </div>
      {error && <Alert tone="danger">{error}</Alert>}
      <Field label="Email">
        <Input type="email" autoComplete="email" required value={email} onChange={(e) => setEmail(e.target.value)} />
      </Field>
      <Field label="Пароль">
        <Input
          type="password"
          autoComplete="current-password"
          required
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
      </Field>
      <Button type="submit" variant="primary" loading={pending}>
        Войти
      </Button>
    </form>
  );
}

// Только относительные пути внутри приложения — без open redirect.
function safeNext(value: string | null): string {
  return value && value.startsWith("/") && !value.startsWith("//") ? value : "/admin";
}
