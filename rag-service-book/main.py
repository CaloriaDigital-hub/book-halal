import os
import logging

from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from openai import OpenAI
from dotenv import load_dotenv
import pytesseract
from graph import build_graph, delete_graph

from vectorstore import search, index_chunks, delete_book_index
from vectorstore import hybrid_search
from chunker import pdf_to_chunks
from ocr import ocr_from_url

load_dotenv()

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Backend URL for downloading book images (use container name in Docker)
BACKEND_URL = os.getenv("BACKEND_URL", "http://book_halal_backend:8090")

# Set Tesseract path
pytesseract.pytesseract.tesseract_cmd = os.getenv("TESSERACT_PATH", "/usr/bin/tesseract")


app = FastAPI(
    title="Book Halal RAG Service",
    description="Retrieval-Augmented Generation for Islamic books",
    version="1.0.0",
)

llm = OpenAI(
    api_key=os.getenv("LLM_API_KEYS"),
    base_url="https://api.deepseek.com",
)

SYSTEM_PROMPT = """Ты помощник по исламским книгам.
Отвечай только на основе предоставленного контекста.
Если ответа нет в контексте — честно скажи что не знаешь.
Указывай страницы источников в ответе."""


# ─── Models ────────────────────────────────────────────────────────────────────

class AskRequest(BaseModel):
    question: str


class IndexPage(BaseModel):
    page_number: int
    image_url: str  # e.g. http://localhost:8090/static/books/{id}/page-001.jpg


class IndexRequest(BaseModel):
    """Index book by passing page image URLs. Python does OCR locally."""
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


# ─── Routes ────────────────────────────────────────────────────────────────────

@app.get("/health", response_model=HealthResponse, tags=["System"])
def health():
    return {"status": "ok"}


@app.post("/books/{book_id}/ask", response_model=AskResponse, tags=["Reader"])
def ask_book(book_id: str, req: AskRequest):
    """Ask a question about an indexed book using RAG + DeepSeek."""
    if not req.question.strip():
        raise HTTPException(status_code=400, detail="question cannot be empty")

    results = hybrid_search(book_id, req.question, top_k=5)

    if not results:
        return AskResponse(
            answer="Книга ещё не проиндексирована или не содержит текста.",
            sources=[],
        )

    context = "\n\n".join(
        f"[Страница {r['metadata']['page']}]: {r['text']}"
        for r in results
    )

    logger.info(f"ask book_id={book_id} question={req.question!r}")

    response = llm.chat.completions.create(
        model="deepseek-chat",
        messages=[
            {"role": "system", "content": SYSTEM_PROMPT},
            {"role": "user", "content": f"Контекст:\n{context}\n\nВопрос: {req.question}"},
        ],
    )

    answer = response.choices[0].message.content

    return AskResponse(
        answer=answer,
        sources=[
            AskSource(page=r["metadata"]["page"], text=r["text"])
            for r in results
        ],
    )


@app.post("/admin/books/{book_id}/index", response_model=IndexResponse, tags=["Admin"])
def index_book(book_id: str, req: IndexRequest):
    """
    Index a book by page image URLs.
    For each page: download image → OCR → store in Qdrant.
    """
    chunks = []
    metadatas = []
    skipped = 0

    if chunks:
        index_chunks(book_id, chunks, metadatas)
    
    # Строим граф
        pages = [m["page"] for m in metadatas]
        build_graph(book_id, chunks, pages)
        logger.info(f"graph built book_id={book_id}")

    for page in req.pages:
        # Replace localhost with backend container hostname
        fixed_url = page.image_url.replace("http://localhost:8090", BACKEND_URL)
        text = ocr_from_url(fixed_url)
        if not text:
            logger.warning(f"book_id={book_id} page={page.page_number}: OCR returned empty, skipping")
            skipped += 1
            continue
        chunks.append(text)
        metadatas.append({"page": page.page_number, "book_id": book_id})

    if chunks:
        index_chunks(book_id, chunks, metadatas)

    logger.info(f"indexed book_id={book_id} indexed={len(chunks)} skipped={skipped}")
    return IndexResponse(indexed=len(chunks), skipped=skipped)


@app.delete("/admin/books/{book_id}/index", tags=["Admin"])
def delete_index(book_id: str):
    """Delete the entire vector index for a book (e.g. before re-upload)."""
    delete_book_index(book_id)
    delete_graph(book_id)  # добавить
    logger.info(f"deleted index for book_id={book_id}")
    return {"deleted": book_id}
    logger.info(f"deleted index for book_id={book_id}")
    return {"deleted": book_id}


@app.get("/admin/books/{book_id}/index", tags=["Admin"])
def inspect_index(book_id: str, limit: int = 10):
    """
    Debug: return stored chunks for a book from Qdrant.
    limit — max number of chunks to return (default 10).
    """
    from vectorstore import client
    collection_name = f"book_{book_id}"

    if not client.collection_exists(collection_name=collection_name):
        return {"book_id": book_id, "total": 0, "chunks": []}

    collection_info = client.get_collection(collection_name=collection_name)
    count = collection_info.points_count

    if count == 0:
        return {"book_id": book_id, "total": 0, "chunks": []}

    results, _ = client.scroll(
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
    """Index a PDF directly by file path (for manual re-indexing via chunker)."""
    chunks, page_numbers = pdf_to_chunks(pdf_path)
    if not chunks:
        raise HTTPException(status_code=422, detail="No text extracted from PDF")

    metadatas = [{"page": p, "book_id": book_id} for p in page_numbers]
    index_chunks(book_id, chunks, metadatas)

    logger.info(f"indexed pdf book_id={book_id} chunks={len(chunks)}")
    return IndexResponse(indexed=len(chunks))
