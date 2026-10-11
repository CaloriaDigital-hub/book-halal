import { formatSom } from "@/shared/lib/format";
import styles from "./PriceTag.module.css";

export function PriceTag({ price, large }: { price: number; large?: boolean }) {
  const free = price === 0;
  return (
    <span className={[styles.price, free && styles.free, large && styles.large].filter(Boolean).join(" ")}>
      {free ? "Бесплатно" : formatSom(price)}
    </span>
  );
}
