# pdf_text.py
"""
Текст страниц напрямую из текстового слоя PDF — без OCR.

Книги, свёрстанные в InDesign и т.п., содержат точный текст. OCR по картинкам
страниц для них только портит его (кириллица путается с латиницей, теряются
знаки транслитерации). OCR нужен лишь для сканов — там слой пустой.
"""
import logging
import os
import re
from collections import Counter

import fitz

logger = logging.getLogger(__name__)

STATIC_ROOT = os.getenv("STATIC_ROOT", "/app/static")

# Меньше этого — считаем, что текстового слоя на странице нет (скан, картинка)
MIN_PAGE_CHARS = 40

# Перенос InDesign: «благ\nг\nгословение» — строка из одного символа-дефиса,
# а следующая строка повторяет последнюю букву предыдущей.
_INDESIGN_HYPHEN = re.compile(r'(\w)\n\S\n\1')
_FIGURE_LINE = re.compile(r'^(Рис|Puc)\.\s*\S+$')


def book_pdf_path(book_id: str) -> str:
    return os.path.join(STATIC_ROOT, "books", book_id, "book.pdf")


def load_pdf_pages(book_id: str) -> dict[int, str]:
    """{номер страницы (с 1): очищенный текст}. Пусто, если PDF нет или он не читается."""
    path = book_pdf_path(book_id)
    if not os.path.exists(path):
        return {}
    try:
        with fitz.open(path) as doc:
            raw = [page.get_text() for page in doc]
    except Exception as e:
        logger.warning(f"book_id={book_id}: failed to read PDF text layer: {e}")
        return {}

    headers = _running_headers(raw)
    return {i: clean_pdf_text(text, headers) for i, text in enumerate(raw, start=1)}


def clean_pdf_text(text: str, headers: set[str] = frozenset()) -> str:
    text = text.replace('\r', '')
    text = _INDESIGN_HYPHEN.sub(r'\1', text)

    lines = []
    for line in text.split('\n'):
        stripped = line.strip()
        if not stripped or stripped.isdigit() or stripped in headers or _FIGURE_LINE.match(stripped):
            continue
        lines.append(stripped)

    text = '\n'.join(lines)
    text = re.sub(r'-\n', '-', text)   # явный дефис в конце строки: «Аль-\nФатиха»
    text = re.sub(r'\n', ' ', text)    # слой PDF не размечает абзацы — склеиваем строки
    text = re.sub(r'\.{4,}', ' ', text)  # точки-заполнители в оглавлении
    text = re.sub(r' {2,}', ' ', text)
    return text.strip()


def _running_headers(pages: list[str]) -> set[str]:
    """Колонтитулы — короткие строки, которые стоят в начале многих страниц."""
    tops = Counter()
    for text in pages:
        head = [l.strip() for l in text.split('\n') if l.strip()][:3]
        tops.update(l for l in set(head) if not l.isdigit() and len(l) < 60)
    return {line for line, n in tops.items() if n >= max(3, len(pages) // 10)}
