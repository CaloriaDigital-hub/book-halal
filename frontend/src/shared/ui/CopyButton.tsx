"use client";

import { useState } from "react";
import { Button } from "./Button";

export function CopyButton({ text, label = "Скопировать" }: { text: string; label?: string }) {
  const [copied, setCopied] = useState(false);

  async function copy() {
    try {
      await navigator.clipboard.writeText(text);
      setCopied(true);
      setTimeout(() => setCopied(false), 1500);
    } catch {
      // Буфер обмена недоступен (http без localhost) — пользователь скопирует вручную
    }
  }

  return (
    <Button size="sm" onClick={copy}>
      {copied ? "Скопировано ✓" : label}
    </Button>
  );
}
