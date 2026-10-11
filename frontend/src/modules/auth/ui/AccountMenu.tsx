"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { Button } from "@/shared/ui";
import { signOut } from "../api/session-api";
import { setSession, useSession } from "../model/use-session";
import styles from "./AccountMenu.module.css";

/** Правая часть шапки магазина: гость — «Войти», покупатель — «Выйти», админ — ещё и ссылка в админку. */
export function AccountMenu({ loginHref, registerHref, adminHref }: { loginHref: string; registerHref: string; adminHref: string }) {
  const state = useSession();
  const router = useRouter();
  const [pending, setPending] = useState(false);

  if (state.status === "loading") return <div className={styles.placeholder} />;

  if (!state.session.authenticated) {
    return (
      <div className={styles.menu}>
        <Link href={loginHref} className={styles.link}>
          Войти
        </Link>
        <Link href={registerHref} className={styles.cta}>
          Регистрация
        </Link>
      </div>
    );
  }

  async function onSignOut() {
    setPending(true);
    try {
      await signOut();
    } finally {
      setSession({ authenticated: false });
      setPending(false);
      router.replace("/");
      router.refresh();
    }
  }

  return (
    <div className={styles.menu}>
      {state.session.role === "admin" && (
        <Link href={adminHref} className={styles.link}>
          Админка
        </Link>
      )}
      <Button variant="ghost" size="sm" loading={pending} onClick={onSignOut}>
        Выйти
      </Button>
    </div>
  );
}
