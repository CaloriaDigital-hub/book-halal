import type { ReactNode } from "react";
import styles from "./Alert.module.css";

type Props = {
  tone?: "danger" | "warning" | "success" | "info";
  title?: string;
  action?: ReactNode;
  children?: ReactNode;
};

export function Alert({ tone = "info", title, action, children }: Props) {
  return (
    <div className={`${styles.alert} ${styles[tone]}`} role={tone === "danger" ? "alert" : "status"}>
      <div className={styles.body}>
        {title && <p className={styles.title}>{title}</p>}
        {children && <div className={styles.text}>{children}</div>}
      </div>
      {action && <div className={styles.action}>{action}</div>}
    </div>
  );
}
