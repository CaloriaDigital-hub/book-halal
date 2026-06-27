import io
import logging
import urllib.request

from PIL import Image
import pytesseract

logger = logging.getLogger(__name__)


import re

def clean_ocr_text(text: str) -> str:
    # Убираем carriage return
    text = text.replace('\r', '')
    # Дефисный перенос слова
    text = re.sub(r'-\n', '', text)
    # Одиночный \n → пробел
    text = re.sub(r'(?<!\n)\n(?!\n)', ' ', text)
    # Множественные пробелы → один
    text = re.sub(r' {2,}', ' ', text)
    # Больше двух \n подряд → два
    text = re.sub(r'\n{3,}', '\n\n', text)
    return text.strip()


def ocr_from_url(image_url: str) -> str:
    """
    Download an image by URL and extract text using Tesseract OCR.

    Supports Arabic, Russian, and English out of the box.
    Returns empty string if download or OCR fails.
    """
    try:
        with urllib.request.urlopen(image_url, timeout=10) as resp:
            img_bytes = resp.read()

        img = Image.open(io.BytesIO(img_bytes))

        # Use available language packs
        raw_text = pytesseract.image_to_string(img, lang="rus+ara+eng")
        
        return clean_ocr_text(raw_text)

    except Exception as e:
        logger.warning(f"OCR failed for {image_url}: {e}")
        return ""
