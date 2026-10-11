"use client";

import type { ReactNode } from "react";
import { Alert } from "./Alert";
import { Button } from "./Button";
import { ErrorBoundary } from "./ErrorBoundary";
import styles from "./Widget.module.css";

type Props = {
  title: string;
  description?: ReactNode;
  actions?: ReactNode;
  children: ReactNode;
};

/**
 * Независимый блок экрана: своя рамка, свой заголовок и свой ErrorBoundary.
 * Если внутри что-то упало — ломается только этот виджет, а не вся страница.
 */
export function Widget({ title, description, actions, children }: Props) {
  return (
    <section className={styles.widget}>
      <header className={styles.header}>
        <div className={styles.heading}>
          <h2 className={styles.title}>{title}</h2>
          {description && <p className={styles.description}>{description}</p>}
        </div>
        {actions && <div className={styles.actions}>{actions}</div>}
      </header>
      <div className={styles.body}>
        <ErrorBoundary
          fallback={(error, reset) => (
            <Alert
              tone="danger"
              title="Блок упал"
              action={
                <Button size="sm" onClick={reset}>
                  Перезапустить
                </Button>
              }
            >
              {error.message}
            </Alert>
          )}
        >
          {children}
        </ErrorBoundary>
      </div>
    </section>
  );
}
