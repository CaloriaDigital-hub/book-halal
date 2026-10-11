import type { ReactNode } from "react";
import styles from "./Grid.module.css";

/** Раскладка виджетов. Колонки схлопываются в одну на узких экранах. */
export function Grid({ columns = 2, children }: { columns?: 1 | 2 | 3; children: ReactNode }) {
  return <div className={`${styles.grid} ${styles[`cols${columns}`]}`}>{children}</div>;
}
