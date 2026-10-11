"use client";

import Link from "next/link";
import { Alert, Spinner, Widget } from "@/shared/ui";
import { useAsync } from "@/shared/lib/use-async";
import { getCatalog } from "../api/books-api";
import styles from "./CatalogStatsWidget.module.css";

export function CatalogStatsWidget({ catalogHref }: { catalogHref: string }) {
  const { data, error } = useAsync("catalog-stats", getCatalog);

  return (
    <Widget title="Книги" actions={<Link href={catalogHref} className={styles.link}>Каталог →</Link>}>
      {error && <Alert tone="danger" title="Бэкенд не отвечает">{error.message}</Alert>}
      {!data && !error && <Spinner label="Загрузка…" />}
      {data && (
        <div className={styles.stats}>
          <div>
            <p className={styles.value}>{data.length}</p>
            <p className={styles.label}>готовых книг</p>
          </div>
          <div>
            <p className={styles.value}>{data.filter((b) => b.price === 0).length}</p>
            <p className={styles.label}>бесплатных</p>
          </div>
        </div>
      )}
    </Widget>
  );
}
