"use client";

import Link from "next/link";
import { useEffect, useRef, useState, type TouchEvent } from "react";
import { ApiError } from "@/shared/api/http";
import { useAsync } from "@/shared/lib/use-async";
import { Spinner } from "@/shared/ui";
import { addBookmark, getBookmarks, removeBookmark } from "../api/reader-api";
import { useReader } from "../model/use-reader";
import { AskPanel } from "./AskPanel";
import { BookmarksPanel } from "./BookmarksPanel";
import styles from "./ReaderView.module.css";

type Props = {
  bookId: string;
  title: string;
  backHref: string;
  buyHref: string;
};

type Drawer = "ask" | "bookmarks" | null;

export function ReaderView({ bookId, title, backHref, buyHref }: Props) {
  const { state, goTo } = useReader(bookId);
  const [drawer, setDrawer] = useState<Drawer>(null);

  if (state.status === "loading") {
    return (
      <div className={styles.center}>
        <Spinner label="Открываем книгу…" />
      </div>
    );
  }
  if (state.status === "error") {
    const forbidden = state.error instanceof ApiError && state.error.status === 403;
    return (
      <div className={styles.center}>
        <div className={styles.notice}>
          <span className={styles.noticeIcon} aria-hidden>
            {forbidden ? "🔒" : "⚠️"}
          </span>
          <h1 className={styles.noticeTitle}>{forbidden ? "Книга ещё не куплена" : "Не удалось открыть книгу"}</h1>
          <p className={styles.noticeText}>
            {forbidden
              ? "Оплатите книгу — после подтверждения она откроется здесь."
              : state.error instanceof Error
                ? state.error.message
                : String(state.error)}
          </p>
          <div className={styles.noticeActions}>
            {forbidden && (
              <Link href={buyHref} className={styles.noticePrimary}>
                Купить книгу
              </Link>
            )}
            <Link href={backHref} className={styles.noticeSecondary}>
              ← К книге
            </Link>
          </div>
        </div>
      </div>
    );
  }

  return (
    <ReaderLayout
      bookId={bookId}
      title={title}
      backHref={backHref}
      pages={state.pages}
      page={state.page}
      goTo={goTo}
      drawer={drawer}
      setDrawer={setDrawer}
    />
  );
}

function ReaderLayout({
  bookId,
  title,
  backHref,
  pages,
  page,
  goTo,
  drawer,
  setDrawer,
}: {
  bookId: string;
  title: string;
  backHref: string;
  pages: { page_number: number; image_url: string }[];
  page: number;
  goTo: (p: number) => void;
  drawer: Drawer;
  setDrawer: (d: Drawer) => void;
}) {
  const total = pages.length;
  const current = pages[page - 1];
  const bookmarks = useAsync(`bookmarks:${bookId}`, () => getBookmarks(bookId));
  const isBookmarked = !!bookmarks.data?.some((b) => b.page_number === page);
  const [pageInput, setPageInput] = useState(String(page));
  const touchX = useRef<number | null>(null);

  // Поле номера страницы следует за листанием
  const [shownPage, setShownPage] = useState(page);
  if (shownPage !== page) {
    setShownPage(page);
    setPageInput(String(page));
  }

  // Клавиатура: ← →, но не когда пользователь печатает вопрос
  useEffect(() => {
    function onKey(e: KeyboardEvent) {
      const tag = (e.target as HTMLElement).tagName;
      if (tag === "INPUT" || tag === "TEXTAREA") return;
      if (e.key === "ArrowRight") goTo(page + 1);
      if (e.key === "ArrowLeft") goTo(page - 1);
    }
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [page, goTo]);

  // Подгружаем соседние страницы заранее — листание без задержек
  useEffect(() => {
    for (const p of [page + 1, page - 1]) {
      const url = pages[p - 1]?.image_url;
      if (url) new Image().src = url;
    }
  }, [page, pages]);

  async function toggleBookmark() {
    if (isBookmarked) await removeBookmark(bookId, page).catch(() => {});
    else await addBookmark(bookId, page).catch(() => {});
    bookmarks.reload();
  }

  function openPage(p: number) {
    goTo(p);
    if (window.matchMedia("(max-width: 900px)").matches) setDrawer(null);
  }

  function onTouchEnd(e: TouchEvent) {
    if (touchX.current === null) return;
    const dx = e.changedTouches[0].clientX - touchX.current;
    touchX.current = null;
    if (Math.abs(dx) > 50) goTo(dx < 0 ? page + 1 : page - 1);
  }

  return (
    <div className={styles.reader}>
      <header className={styles.topbar}>
        <Link href={backHref} className={styles.back} aria-label="К книге">
          ←
        </Link>
        <span className={styles.title}>{title}</span>
        <div className={styles.tools}>
          <button
            type="button"
            className={isBookmarked ? `${styles.tool} ${styles.toolActive}` : styles.tool}
            onClick={toggleBookmark}
            title={isBookmarked ? "Убрать закладку" : "Добавить закладку"}
          >
            🔖
          </button>
          <button
            type="button"
            className={drawer === "bookmarks" ? `${styles.tool} ${styles.toolActive}` : styles.tool}
            onClick={() => setDrawer(drawer === "bookmarks" ? null : "bookmarks")}
            title="Закладки"
          >
            <span aria-hidden>📑</span>
            <span className={styles.wide}> Закладки</span>
          </button>
          <button
            type="button"
            className={drawer === "ask" ? `${styles.askButton} ${styles.toolActive}` : styles.askButton}
            onClick={() => setDrawer(drawer === "ask" ? null : "ask")}
          >
            ✨ Спросить ИИ
          </button>
        </div>
      </header>

      <div className={styles.body}>
        <div
          className={styles.stage}
          onTouchStart={(e) => (touchX.current = e.touches[0].clientX)}
          onTouchEnd={onTouchEnd}
        >
          <button type="button" className={`${styles.nav} ${styles.prev}`} onClick={() => goTo(page - 1)} disabled={page <= 1} aria-label="Предыдущая страница">
            ‹
          </button>
          {current && (
            // eslint-disable-next-line @next/next/no-img-element -- страницы — картинки из /static Go-бэка
            <img key={current.image_url} src={current.image_url} alt={`Страница ${page}`} className={styles.page} />
          )}
          <button type="button" className={`${styles.nav} ${styles.next}`} onClick={() => goTo(page + 1)} disabled={page >= total} aria-label="Следующая страница">
            ›
          </button>
        </div>

        {drawer && (
          <aside className={styles.drawer}>
            <div className={styles.drawerHead}>
              <span>{drawer === "ask" ? "Вопрос по книге" : "Закладки"}</span>
              <button type="button" className={styles.close} onClick={() => setDrawer(null)} aria-label="Закрыть">
                ✕
              </button>
            </div>
            <div className={styles.drawerBody}>
              {/* Панель вопросов не размонтируется при переключении, чтобы не терять историю */}
              <div hidden={drawer !== "ask"} className={styles.fill}>
                <AskPanel bookId={bookId} onOpenPage={openPage} />
              </div>
              {drawer === "bookmarks" && (
                <BookmarksPanel bookmarks={bookmarks.data} error={bookmarks.error} currentPage={page} onOpenPage={openPage} />
              )}
            </div>
          </aside>
        )}
      </div>

      <footer className={styles.bottombar}>
        <input
          type="range"
          min={1}
          max={total}
          value={page}
          onChange={(e) => goTo(Number(e.target.value))}
          className={styles.slider}
          aria-label="Страница"
        />
        <form
          className={styles.pageForm}
          onSubmit={(e) => {
            e.preventDefault();
            goTo(Number(pageInput));
          }}
        >
          <input
            className={styles.pageInput}
            inputMode="numeric"
            value={pageInput}
            onChange={(e) => setPageInput(e.target.value.replace(/\D/g, ""))}
            onBlur={() => goTo(Number(pageInput))}
            aria-label="Номер страницы"
          />
          <span className={styles.total}>из {total}</span>
        </form>
      </footer>
    </div>
  );
}
