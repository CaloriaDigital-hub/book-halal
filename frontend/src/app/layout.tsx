import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Book Halal — админка",
  description: "Управление книгами и RAG-индексом",
};

export default function RootLayout({ children }: LayoutProps<"/">) {
  return (
    <html lang="ru">
      <body>{children}</body>
    </html>
  );
}
