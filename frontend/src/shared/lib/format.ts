const dateTime = new Intl.DateTimeFormat("ru-RU", {
  day: "2-digit",
  month: "2-digit",
  hour: "2-digit",
  minute: "2-digit",
  timeZone: "Asia/Bishkek",
});

/** 1500 → «1 500 сом». Цены в API — в сомах. */
export function formatSom(amount: number): string {
  return `${amount.toLocaleString("ru-RU")} сом`;
}

/** Дата и время по Бишкеку: «11.10, 22:55». */
export function formatDateTime(iso: string): string {
  return dateTime.format(new Date(iso));
}
