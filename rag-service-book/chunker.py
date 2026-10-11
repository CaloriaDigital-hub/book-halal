import fitz
import re

def pdf_to_chunks(
    pdf_path: str,
    chunk_size: int = 500,
    overlap: int = 100
) -> tuple[list[str], list[int]]:
    doc = fitz.open(pdf_path)
    chunks = []
    page_numbers = []

    for page_num, page in enumerate(doc, start=1):
        text = page.get_text().strip()
        if not text:
            continue

        # Чистим мусор
        text = re.sub(r'\n+', '\n', text)          # множественные переносы
        text = re.sub(r'[ \t]+', ' ', text)         # множественные пробелы
        text = re.sub(r'-\n(\w)', r'\1', text)      # переносы слов: "ис-\nлам" → "ислам"

        # Разбиваем на предложения, не на символы
        sentences = re.split(r'(?<=[.!?])\s+', text)

        current_chunk = []
        current_len = 0

        for sentence in sentences:
            sentence_len = len(sentence)

            if current_len + sentence_len > chunk_size and current_chunk:
                chunk_text = ' '.join(current_chunk).strip()
                if chunk_text:
                    chunks.append(chunk_text)
                    page_numbers.append(page_num)

                # Overlap — берём последние предложения
                overlap_chunk = []
                overlap_len = 0
                for s in reversed(current_chunk):
                    if overlap_len + len(s) <= overlap:
                        overlap_chunk.insert(0, s)
                        overlap_len += len(s)
                    else:
                        break

                current_chunk = overlap_chunk
                current_len = overlap_len

            current_chunk.append(sentence)
            current_len += sentence_len

        # Последний чанк страницы
        if current_chunk:
            chunk_text = ' '.join(current_chunk).strip()
            if chunk_text:
                chunks.append(chunk_text)
                page_numbers.append(page_num)

    doc.close()
    return chunks, page_numbers

# ─── Нарезка OCR-текста страницы на чанки ──────────────────────

_SENTENCE_END = re.compile(r'(?<=[.!?…])\s+')


def split_text(
    text: str,
    count_tokens,
    max_tokens: int = 300,
    overlap_tokens: int = 60,
) -> list[str]:
    """
    Режет текст страницы на чанки по предложениям, каждый не длиннее max_tokens
    токенов эмбеддера (e5 и реранкер видят только первые 512). Соседние чанки
    перекрываются хвостом предыдущего примерно на overlap_tokens.
    """
    units: list[tuple[str, int]] = []
    for paragraph in re.split(r'\n{2,}', text):
        for sentence in _SENTENCE_END.split(paragraph.strip()):
            sentence = sentence.strip()
            if not sentence:
                continue
            n = count_tokens(sentence)
            if n <= max_tokens:
                units.append((sentence, n))
            else:
                units.extend(_split_long(sentence, count_tokens, max_tokens))

    chunks: list[str] = []
    current: list[tuple[str, int]] = []
    current_len = 0

    for unit in units:
        if current and current_len + unit[1] > max_tokens:
            chunks.append(' '.join(t for t, _ in current))

            # Overlap — последние предложения предыдущего чанка
            overlap: list[tuple[str, int]] = []
            overlap_len = 0
            for t, n in reversed(current):
                if overlap_len + n > overlap_tokens:
                    break
                overlap.insert(0, (t, n))
                overlap_len += n
            current, current_len = overlap, overlap_len

            while current and current_len + unit[1] > max_tokens:
                current_len -= current.pop(0)[1]

        current.append(unit)
        current_len += unit[1]

    if current:
        chunks.append(' '.join(t for t, _ in current))
    return chunks


def _split_long(sentence: str, count_tokens, max_tokens: int) -> list[tuple[str, int]]:
    """Предложение длиннее лимита (частый случай при плохом OCR) режем по словам."""
    pieces: list[tuple[str, int]] = []
    words: list[str] = []
    words_len = 0
    for word in sentence.split():
        n = count_tokens(word)
        if words and words_len + n > max_tokens:
            pieces.append((' '.join(words), words_len))
            words, words_len = [], 0
        words.append(word)
        words_len += n
    if words:
        pieces.append((' '.join(words), words_len))
    return pieces
