// Ответы Go-бэка — на английском; покупателю показываем понятный текст.
const KNOWN: [RegExp, string][] = [
  [/receipt must be/, "Чек должен быть картинкой (jpg, png, webp) или PDF"],
  [/larger than 10 MB|invalid form data/, "Файл слишком большой — максимум 10 МБ"],
  [/receipt file is required/, "Прикрепите чек"],
  [/expected phone number/, "Укажите номер телефона, например +996 555 123 456"],
  [/expected telegram/, "Укажите Telegram в виде @username или номер телефона"],
  [/already pending or approved/, "Чек по этой книге уже отправлен — дождитесь проверки"],
  [/book is free/, "Эта книга бесплатная — оплата не нужна"],
];

export function humanize(error: unknown): string {
  const message = error instanceof Error ? error.message : String(error);
  return KNOWN.find(([re]) => re.test(message))?.[1] ?? message;
}
