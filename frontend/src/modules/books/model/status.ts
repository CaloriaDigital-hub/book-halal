import { formatSom } from "@/shared/lib/format";
import type { Tone } from "@/shared/ui";

const STATUS: Record<string, { label: string; tone: Tone }> = {
  ready: { label: "Готова", tone: "success" },
  processing: { label: "Обработка", tone: "warning" },
  error: { label: "Ошибка", tone: "danger" },
};

export function bookStatus(status: string) {
  return STATUS[status] ?? { label: status || "—", tone: "neutral" as Tone };
}

export function formatPrice(price: number) {
  return price === 0 ? "Бесплатно" : formatSom(price);
}
