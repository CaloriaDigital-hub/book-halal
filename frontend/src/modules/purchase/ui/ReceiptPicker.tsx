"use client";

import { useEffect, useRef, useState, type DragEvent } from "react";
import styles from "./ReceiptPicker.module.css";

const ACCEPT = "image/jpeg,image/png,image/webp,application/pdf";
const MAX_BYTES = 10 * 1024 * 1024;

/** Выбор чека: клик или перетаскивание, превью картинки. */
export function ReceiptPicker({ file, onChange, onError }: { file: File | null; onChange: (f: File | null) => void; onError: (m: string) => void }) {
  const [dragging, setDragging] = useState(false);
  const [preview, setPreview] = useState<string | null>(null);
  const previewRef = useRef<string | null>(null);

  // Освобождаем object URL превью при уходе со страницы
  useEffect(() => () => {
    if (previewRef.current) URL.revokeObjectURL(previewRef.current);
  }, []);

  function pick(f: File | undefined) {
    if (!f) return;
    if (!ACCEPT.split(",").includes(f.type)) return onError("Нужна картинка (jpg, png, webp) или PDF");
    if (f.size > MAX_BYTES) return onError("Файл больше 10 МБ — сожмите скриншот или пришлите PDF");
    if (previewRef.current) URL.revokeObjectURL(previewRef.current);
    previewRef.current = f.type.startsWith("image/") ? URL.createObjectURL(f) : null;
    setPreview(previewRef.current);
    onChange(f);
  }

  function onDrop(e: DragEvent) {
    e.preventDefault();
    setDragging(false);
    pick(e.dataTransfer.files[0]);
  }

  return (
    <label
      className={[styles.zone, dragging && styles.dragging, file && styles.filled].filter(Boolean).join(" ")}
      onDragOver={(e) => {
        e.preventDefault();
        setDragging(true);
      }}
      onDragLeave={() => setDragging(false)}
      onDrop={onDrop}
    >
      <input type="file" accept={ACCEPT} className={styles.input} onChange={(e) => pick(e.target.files?.[0])} />
      {file ? (
        <div className={styles.selected}>
          {preview ? (
            // eslint-disable-next-line @next/next/no-img-element -- локальное превью выбранного файла
            <img src={preview} alt="Чек" className={styles.preview} />
          ) : (
            <span className={styles.pdf}>PDF</span>
          )}
          <div className={styles.fileInfo}>
            <span className={styles.fileName}>{file.name}</span>
            <span className={styles.change}>Нажмите, чтобы выбрать другой файл</span>
          </div>
        </div>
      ) : (
        <div className={styles.empty}>
          <span className={styles.icon} aria-hidden>
            🧾
          </span>
          <span className={styles.label}>Прикрепите чек</span>
          <span className={styles.sub}>Скриншот из банковского приложения или PDF · до 10 МБ</span>
        </div>
      )}
    </label>
  );
}
