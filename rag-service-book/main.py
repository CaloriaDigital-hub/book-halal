# main.py
import os
import logging

from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import pytesseract
from dotenv import load_dotenv

# Новые компоненты
from embeddings import SentenceTransformerEmbedder
from stores import QdrantVectorStore
from retrievers import HybridRetriever
from generators import DeepSeekGenerator
from pipeline import BookRAG

from chunker import pdf_to_chunks
from ocr import ocr_from_url

load_dotenv()

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

BACKEND_URL = os.getenv("BACKEND_URL", "http://book_halal_backend:8090")
pytesseract.pytesseract.tesseract_cmd = os.getenv("TESSERACT_PATH", "/usr/bin/tesseract")

# ═══════════════════════════════════════════════════════════
# Инициализация RAG
# ═══════════════════════════════════════════════════════════
embedder = SentenceTransformerEmbedder()
vector_store = QdrantVectorStore()
retriever = HybridRetriever(embedder=embedder, vector_store=vector_store)
generator = DeepSeekGenerator()
rag = BookRAG(embedder=embedder, vector_store=vector_store,
              retriever=retriever, generator=generator)

app = FastAPI(
    title="Book Halal RAG Service",
    description="Retrieval-Augmented Generation for Islamic books",
    version="2.0.0",
)

# ─── Модели ─────────────────────────────────────────────

class AskRequest(BaseModel):
    question: str

class IndexPage(BaseModel):
    page_number: int
    image_url: str

class IndexRequest(BaseModel):
    pages: list[IndexPage]

class AskSource(BaseModel):
    page: int
    text: str

class AskResponse(BaseModel):
    answer: str
    sources: list[AskSource]

class IndexResponse(BaseModel):
    indexed: int
    skipped: int = 0

class HealthResponse(BaseModel):
    status: str


# ─── Роуты ──────────────────────────────────────────────

@app.get("/health", response_model=HealthResponse, tags=["System"])
def health():
    return {"status": "ok"}


@app.post("/books/{book_id}/ask", response_model=AskResponse, tags=["Reader"])
def ask_book(book_id: str, req: AskRequest):
    if not req.question.strip():
        raise HTTPException(status_code=400, detail="question cannot be empty")

    result = rag.ask(book_id, req.question)
    return AskResponse(
        answer=result["answer"],
        sources=[AskSource(**s) for s in result["sources"]],
    )


@app.post("/admin/books/{book_id}/index", response_model=IndexResponse, tags=["Admin"])
def index_book(book_id: str, req: IndexRequest):
    chunks = []
    metadatas = []
    skipped = 0

    for page in req.pages:
        fixed_url = page.image_url.replace("http://localhost:8090", BACKEND_URL)
        text = ocr_from_url(fixed_url)
        if not text:
            logger.warning(f"book_id={book_id} page={page.page_number}: OCR empty, skipping")
            skipped += 1
            continue
        chunks.append(text)
        metadatas.append({"page": page.page_number, "book_id": book_id})

    if chunks:
        rag.index(book_id, chunks, metadatas)

    logger.info(f"indexed book_id={book_id} indexed={len(chunks)} skipped={skipped}")
    return IndexResponse(indexed=len(chunks), skipped=skipped)


@app.delete("/admin/books/{book_id}/index", tags=["Admin"])
def delete_index(book_id: str):
    rag.delete_index(book_id)
    logger.info(f"deleted index for book_id={book_id}")
    return {"deleted": book_id}


@app.get("/admin/books/{book_id}/index", tags=["Admin"])
def inspect_index(book_id: str, limit: int = 10):
    collection_name = f"book_{book_id}"

    # Используем клиент Qdrant напрямую для скролла (упрощённо)
    qdrant_client = rag.vector_store.client
    if not qdrant_client.collection_exists(collection_name=collection_name):
        return {"book_id": book_id, "total": 0, "chunks": []}

    collection_info = qdrant_client.get_collection(collection_name=collection_name)
    count = collection_info.points_count

    if count == 0:
        return {"book_id": book_id, "total": 0, "chunks": []}

    results, _ = qdrant_client.scroll(
        collection_name=collection_name,
        limit=limit,
        with_payload=True,
        with_vectors=False,
    )

    chunks = [
        {"page": point.payload.get("page"), "text": point.payload.get("text")}
        for point in results
    ]

    return {"book_id": book_id, "total": count, "chunks": chunks}


@app.post("/admin/books/{book_id}/index-pdf", response_model=IndexResponse, tags=["Admin"])
def index_pdf(book_id: str, pdf_path: str):
    chunks, page_numbers = pdf_to_chunks(pdf_path)
    if not chunks:
        raise HTTPException(status_code=422, detail="No text extracted from PDF")

    metadatas = [{"page": p, "book_id": book_id} for p in page_numbers]
    rag.index(book_id, chunks, metadatas)

    logger.info(f"indexed pdf book_id={book_id} chunks={len(chunks)}")
    return IndexResponse(indexed=len(chunks))


@app.get("/debug/page/{book_id}")
def debug_page(book_id: str, page: int):
    collection_name = f"book_{book_id}"
    points, _ = rag.vector_store.client.scroll(
        collection_name=collection_name,
        scroll_filter={
            "must": [
                {"key": "book_id", "match": {"value": book_id}},
                {"key": "page", "match": {"value": page}}
            ]
        },
        limit=1,
    )
    if points:
        payload = points[0].payload
        return {"found": True, "page": payload.get("page"), "text": payload.get("text", "")}
    return {"found": False}