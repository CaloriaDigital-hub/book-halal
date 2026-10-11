import type { ContactType } from "./types";

const LABEL: Record<ContactType, string> = { phone: "Телефон", whatsapp: "WhatsApp", telegram: "Telegram" };

/** Подпись и ссылка, по которой можно сразу написать покупателю. */
export function contactLink(type: ContactType, value: string): { label: string; href: string | null } {
  const digits = value.replace(/\D/g, "");
  switch (type) {
    case "telegram":
      return { label: LABEL[type], href: value.startsWith("@") ? `https://t.me/${value.slice(1)}` : null };
    case "whatsapp":
      return { label: LABEL[type], href: `https://wa.me/${digits}` };
    default:
      return { label: LABEL[type], href: `tel:${value.startsWith("+") ? "+" : ""}${digits}` };
  }
}
