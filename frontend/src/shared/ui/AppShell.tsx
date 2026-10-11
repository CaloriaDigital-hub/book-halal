import { Suspense, type ReactNode } from "react";
import { ActiveNav, NavLinks, type NavItem } from "./AppShellNav";
import styles from "./AppShell.module.css";

export type { NavItem };

type Props = {
  brand: ReactNode;
  nav: NavItem[];
  footer?: ReactNode;
  children: ReactNode;
};

export function AppShell({ brand, nav, footer, children }: Props) {
  return (
    <div className={styles.shell}>
      <aside className={styles.sidebar}>
        <div className={styles.brand}>{brand}</div>
        {/* Подсветка активного пункта зависит от URL — догружается отдельно, каркас рендерится сразу. */}
        <Suspense fallback={<NavLinks nav={nav} pathname={null} />}>
          <ActiveNav nav={nav} />
        </Suspense>
        {footer && <div className={styles.footer}>{footer}</div>}
      </aside>
      <main className={styles.main}>{children}</main>
    </div>
  );
}
