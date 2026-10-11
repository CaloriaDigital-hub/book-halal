import Link from "next/link";
import type { ReactNode } from "react";
import styles from "./StoreShell.module.css";

type Props = {
  brand: string;
  homeHref: string;
  actions?: ReactNode;
  children: ReactNode;
};

/** Каркас витрины: шапка с логотипом и действиями, контент по центру, подвал. */
export function StoreShell({ brand, homeHref, actions, children }: Props) {
  return (
    <div className={styles.shell}>
      <header className={styles.header}>
        <div className={styles.headerInner}>
          <Link href={homeHref} className={styles.brand}>
            <span className={styles.logo} aria-hidden>
              📖
            </span>
            {brand}
          </Link>
          {actions}
        </div>
      </header>
      <main className={styles.main}>{children}</main>
      <footer className={styles.footer}>© {brand}</footer>
    </div>
  );
}
