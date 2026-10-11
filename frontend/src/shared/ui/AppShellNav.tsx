"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import styles from "./AppShellNav.module.css";

export type NavItem = { href: string; label: string; exact?: boolean };

export function ActiveNav({ nav }: { nav: NavItem[] }) {
  return <NavLinks nav={nav} pathname={usePathname()} />;
}

export function NavLinks({ nav, pathname }: { nav: NavItem[]; pathname: string | null }) {
  const isActive = (item: NavItem) =>
    pathname !== null &&
    (item.exact ? pathname === item.href : pathname === item.href || pathname.startsWith(`${item.href}/`));

  return (
    <nav className={styles.nav}>
      {nav.map((item) => (
        <Link key={item.href} href={item.href} className={isActive(item) ? `${styles.link} ${styles.active}` : styles.link}>
          {item.label}
        </Link>
      ))}
    </nav>
  );
}
