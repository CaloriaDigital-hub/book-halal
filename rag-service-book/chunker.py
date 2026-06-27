import fitz  # pymupdf


def pdf_to_chunks(pdf_path: str, chunk_size: int = 500) -> tuple[list[str], list[int]]:
    """
    Extract text chunks from a PDF file.

    Returns:
        chunks: list of text strings
        page_numbers: list of page numbers (1-based) corresponding to each chunk
    """
    doc = fitz.open(pdf_path)
    chunks = []
    page_numbers = []

    for page_num, page in enumerate(doc, start=1):
        text = page.get_text().strip()
        if not text:
            continue

        # Split page text into fixed-size chunks with overlap
        for i in range(0, len(text), chunk_size):
            chunk = text[i : i + chunk_size].strip()
            if chunk:
                chunks.append(chunk)
                page_numbers.append(page_num)

    doc.close()
    return chunks, page_numbers