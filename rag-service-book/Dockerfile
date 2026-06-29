FROM python:3.11-slim

WORKDIR /app


RUN apt-get update && apt-get install -y \
    tesseract-ocr \
    tesseract-ocr-rus \
    tesseract-ocr-ara \
    && rm -rf /var/lib/apt/lists/*

# Install dependencies first (layer caching)
COPY requirements.txt .
RUN pip install uv && \
    uv pip install --system -r requirements.txt

# Copy source code
COPY . .

# Persist ChromaDB data via volume
VOLUME ["/app/chroma_db"]

EXPOSE 8000

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
