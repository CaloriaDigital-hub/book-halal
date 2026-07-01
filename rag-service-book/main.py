import os
import logging

from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from openai import OpenAI
from dotenv import load_dotenv
import pytesseract
from graph_legacy import build_graph, delete_graph

from vectorstore_legacy import (
    search,              # не используется, но пусть будет
    index_chunks,
    delete_book_index,
    hybrid_search,
    client as qdrant_client,   # <- для inspect_index
)
from chunker import pdf_to_chunks
from ocr import ocr_from_url

load_dotenv()

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

BACKEND_URL = os.getenv("BACKEND_URL", "http://book_halal_backend:8090")
pytesseract.pytesseract.tesseract_cmd = os.getenv("TESSERACT_PATH", "/usr/bin/tesseract")

app = FastAPI(
    title="Book Halal RAG Service",
    description="Retrieval-Augmented Generation for Islamic books",
    version="2.0.0",
)

llm = OpenAI(
    api_key=os.getenv("LLM_API_KEYS"),
    base_url="https://api.deepseek.com",
)

SYSTEM_PROMPT = """Ты — ассистент, который отвечает на вопросы СТРОГО на основе контекста.
Контекст взят из исламской книги.

Правила:
1. Используй ТОЛЬКО информацию, приведённую в контексте.
2. Если в контексте нет точного ответа, напиши ровно: "В данной книге ответ не найден."
3. Отвечай на том же языке, что и вопрос.
4. В ответе обязательно указывай источник в формате [Источник: страница X].
5. Не добавляй никаких внешних знаний, даже если ты знаешь ответ.
6. Если контекст содержит аяты или хадисы, цитируй их точно.
7. Не придумывай связи между фрагментами, если они не очевидны."""


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

    results = hybrid_search(book_id, req.question, top_k=10)

    if not results:
        return AskResponse(
            answer="Книга ещё не проиндексирована или не содержит текста.",
            sources=[],
        )

    context_parts = []
    for i, r in enumerate(results, 1):
        page = r["metadata"]["page"]
        text = r["text"]
        context_parts.append(f"[Фрагмент {i}]\n{text}\n[Источник: страница {page}]\n")
    context = "\n".join(context_parts)

    logger.info(f"ask book_id={book_id} question={req.question!r}")

    response = llm.chat.completions.create(
        model="deepseek-chat",
        temperature=0.1,
        messages=[
            {"role": "system", "content": SYSTEM_PROMPT},
            {"role": "user", "content": f"Контекст:\n---\n{context}\n---\n\nВопрос: {req.question}"},
        ],
    )

    answer = response.choices[0].message.content

    return AskResponse(
        answer=answer,
        sources=[
            AskSource(page=r["metadata"]["page"], text=r["text"])
            for r in results[:5]
        ],
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
        index_chunks(book_id, chunks, metadatas)
        # Строим граф для улучшенного поиска
        pages = [m["page"] for m in metadatas]
        build_graph(book_id, chunks, pages)
        logger.info(f"graph built book_id={book_id}")

    logger.info(f"indexed book_id={book_id} indexed={len(chunks)} skipped={skipped}")
    return IndexResponse(indexed=len(chunks), skipped=skipped)


@app.delete("/admin/books/{book_id}/index", tags=["Admin"])
def delete_index(book_id: str):
    delete_book_index(book_id)
    delete_graph(book_id)
    logger.info(f"deleted index for book_id={book_id}")
    return {"deleted": book_id}


@app.get("/admin/books/{book_id}/index", tags=["Admin"])
def inspect_index(book_id: str, limit: int = 10):
    collection_name = f"book_{book_id}"

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
    index_chunks(book_id, chunks, metadatas)
    # Строим граф и для PDF
    build_graph(book_id, chunks, page_numbers)
    logger.info(f"graph built for pdf book_id={book_id}")

    logger.info(f"indexed pdf book_id={book_id} chunks={len(chunks)}")
    return IndexResponse(indexed=len(chunks))