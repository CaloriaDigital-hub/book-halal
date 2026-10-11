import styles from "./Cover.module.css";

/** Обложка книги с заглушкой, если картинки нет. */
export function Cover({ src, title, className }: { src: string; title: string; className?: string }) {
  return (
    <div className={[styles.cover, className].filter(Boolean).join(" ")}>
      {src ? (
        // eslint-disable-next-line @next/next/no-img-element -- обложки отдаёт Go-бэк из /static
        <img src={src} alt={title} loading="lazy" />
      ) : (
        <span className={styles.placeholder}>{title}</span>
      )}
    </div>
  );
}
