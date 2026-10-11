/** Только относительные пути внутри приложения — без open redirect. */
export function safeNext(value: string | null): string | null {
  return value && value.startsWith("/") && !value.startsWith("//") ? value : null;
}
