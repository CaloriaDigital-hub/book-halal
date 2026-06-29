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