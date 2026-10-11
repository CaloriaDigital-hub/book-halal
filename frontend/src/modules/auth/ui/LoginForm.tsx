"use client";

import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { useState, type FormEvent } from "react";
import { Alert, Button, Field, Input } from "@/shared/ui";
import { signIn } from "../api/session-api";
import { safeNext } from "../model/safe-next";
import { setSession } from "../model/use-session";
import styles from "./AuthCard.module.css";

export function LoginForm({ registerHref }: { registerHref: string }) {
  const router = useRouter();
  const params = useSearchParams();
  const next = safeNext(params.get("next"));
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function onSubmit(event: FormEvent) {
    event.preventDefault();
    setPending(true);
    setError(null);
    try {
      const { role } = await signIn(email, password);
      setSession({ authenticated: true, role });
      router.replace(next ?? (role === "admin" ? "/admin" : "/"));
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
      setPending(false);
    }
  }

  const registerLink = next ? `${registerHref}?next=${encodeURIComponent(next)}` : registerHref;

  return (
    <form className={styles.card} onSubmit={onSubmit}>
      <div>
        <h1 className={styles.title}>Вход</h1>
        <p className={styles.subtitle}>Чтобы покупать и читать книги</p>
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
      <p className={styles.switch}>
        Нет аккаунта? <Link href={registerLink}>Зарегистрироваться</Link>
      </p>
    </form>
  );
}
