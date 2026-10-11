import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Book Halal — исламские книги онлайн",
  description: "Покупайте и читайте исламские книги, задавайте вопросы по тексту",
};

export default function RootLayout({ children }: LayoutProps<"/">) {
  return (
    <html lang="ru">
      <body>{children}</body>
    </html>
  );
}
