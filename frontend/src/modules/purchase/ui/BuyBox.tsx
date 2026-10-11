"use client";

import Link from "next/link";
import { useEffect } from "react";
import { formatSom } from "@/shared/lib/format";
import { Alert, Spinner } from "@/shared/ui";
import { useMyPurchase } from "../model/use-my-purchase";
import styles from "./BuyBox.module.css";

type Props = {
  bookId: string;
  price: number;
  hasAccess: boolean;
  /** null — сессия ещё загружается */
  signedIn: boolean | null;
  readHref: string;
  buyHref: string;
  loginHref: string;
  /** Оплату подтвердили, пока страница открыта — экран перечитает книгу (has_access). */
  onAccessGranted: () => void;
};

/** Главная кнопка на странице книги — зависит от того, что пользователю доступно сейчас. */
export function BuyBox({ bookId, price, hasAccess, signedIn, readHref, buyHref, loginHref, onAccessGranted }: Props) {
  const needsPurchase = !hasAccess && price > 0 && signedIn === true;
  const { data: purchase, loading } = useMyPurchase(bookId, needsPurchase);
  const approved = purchase?.status === "approved";

  useEffect(() => {
    if (approved && !hasAccess) onAccessGranted();
  }, [approved, hasAccess, onAccessGranted]);

  if (hasAccess) {
    return (
      <div className={styles.box}>
        <Link href={readHref} className={styles.primary}>
          Читать книгу
        </Link>
        <p className={styles.hint}>{price > 0 ? "Книга куплена — доступна навсегда" : "Бесплатная книга"}</p>
      </div>
    );
  }
  if (signedIn === null || (needsPurchase && loading && purchase === undefined)) {
    return <Spinner />;
  }
  if (!signedIn) {
    return (
      <div className={styles.box}>
        <Link href={loginHref} className={styles.primary}>
          Войти, чтобы {price > 0 ? "купить" : "читать"}
        </Link>
        <p className={styles.hint}>Регистрация — по email за минуту</p>
      </div>
    );
  }

  if (purchase?.status === "pending") {
    return (
      <Alert tone="info" title="Чек на проверке ⏳">
        Мы получили ваш чек и проверяем оплату. Книга откроется здесь автоматически — страницу можно не обновлять.
      </Alert>
    );
  }

  return (
    <div className={styles.box}>
      {purchase?.status === "rejected" && (
        <Alert tone="warning" title="Оплата не подтверждена">
          Мы не смогли найти этот платёж. Проверьте сумму и реквизиты и отправьте чек ещё раз — или дождитесь, мы свяжемся с вами.
        </Alert>
      )}
      <Link href={buyHref} className={styles.primary}>
        {purchase?.status === "rejected" ? "Отправить чек заново" : `Купить за ${formatSom(price)}`}
      </Link>
      <p className={styles.hint}>Перевод по QR или реквизитам, затем — фото чека</p>
    </div>
  );
}
