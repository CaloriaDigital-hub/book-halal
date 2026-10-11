import { SignOutButton } from "@/modules/auth";
import { AppShell, type NavItem } from "@/shared/ui";

const NAV: NavItem[] = [
  { href: "/admin", label: "Обзор", exact: true },
  { href: "/admin/books", label: "Книги" },
];

export default function AdminLayout({ children }: LayoutProps<"/admin">) {
  return (
    <AppShell brand="Book Halal" nav={NAV} footer={<SignOutButton />}>
      {children}
    </AppShell>
  );
}
