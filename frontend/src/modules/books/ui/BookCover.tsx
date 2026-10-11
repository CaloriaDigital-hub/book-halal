import styles from "./BookCover.module.css";

export function BookCover({ src, title, size = "sm" }: { src: string; title: string; size?: "sm" | "lg" }) {
  return (
    <div className={`${styles.cover} ${styles[size]}`}>
      {src ? (
        // eslint-disable-next-line @next/next/no-img-element -- обложки отдаёт Go-бэк, оптимизатор Next тут не нужен
        <img src={src} alt={title} loading="lazy" />
      ) : (
        <span className={styles.placeholder}>{title.slice(0, 1).toUpperCase()}</span>
      )}
    </div>
  );
}
