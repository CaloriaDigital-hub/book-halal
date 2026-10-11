import Link from "next/link";
import type { ReactNode } from "react";
import styles from "./PageHeader.module.css";

type Props = {
  title: ReactNode;
  subtitle?: ReactNode;
  back?: { href: string; label: string };
};

export function PageHeader({ title, subtitle, back }: Props) {
  return (
    <header className={styles.header}>
      {back && (
        <Link href={back.href} className={styles.back}>
          ← {back.label}
        </Link>
      )}
      <h1 className={styles.title}>{title}</h1>
      {subtitle && <p className={styles.subtitle}>{subtitle}</p>}
    </header>
  );
}
