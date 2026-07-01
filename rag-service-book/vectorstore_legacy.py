import uuid
import os
from qdrant_client import QdrantClient
from qdrant_client.models import Distance, VectorParams, PointStruct
from sentence_transformers import SentenceTransformer, CrossEncoder

from graph_legacy import load_graph, search_graph

# ═══════════════════════════════════════════════════════════════
# Модели
# ═══════════════════════════════════════════════════════════════

# Эмбеддер – multilingual-e5-large (1024‑мерные векторы)
embed_model = SentenceTransformer("intfloat/multilingual-e5-large")

# Кросс‑энкодер для реранкинга (multilingual, понимает русский и арабский)
reranker = CrossEncoder("BAAI/bge-reranker-v2-m3", max_length=512)

# Клиент Qdrant
client = QdrantClient(url=os.getenv("QDRANT_URL", "http://localhost:6333"))


# ═══════════════════════════════════════════════════════════════
# Вспомогательные функции
# ═══════════════════════════════════════════════════════════════

def _embed_documents(texts: list[str]) -> list[list[float]]:
    """Эмбеддинги для индексации (с префиксом 'passage: ')."""
    return embed_model.encode(
        [f"passage: {t}" for t in texts],
        normalize_embeddings=True,
    ).tolist()


def _embed_query(query: str) -> list[float]:
    """Эмбеддинг запроса (с префиксом 'query: ')."""
    return embed_model.encode(
        f"query: {query}",
        normalize_embeddings=True,
    ).tolist()


def rerank(query: str, documents: list[dict], top_k: int) -> list[dict]:
    """
    Ранжирует список документов (словарей с ключом 'text') 
    и возвращает top_k наиболее релевантных.
    """
    if not documents:
        return []

    texts = [doc["text"] for doc in documents]
    pairs = [[query, t] for t in texts]
    scores = reranker.predict(pairs)  # список float

    # Сортируем по убыванию скоров
    scored = list(zip(documents, scores))
    scored.sort(key=lambda x: x[1], reverse=True)
    return [doc for doc, _ in scored[:top_k]]


# ═══════════════════════════════════════════════════════════════
# Управление коллекциями и индексация
# ═══════════════════════════════════════════════════════════════

def get_or_create_collection(book_id: str) -> str:
    collection_name = f"book_{book_id}"
    if not client.collection_exists(collection_name=collection_name):
        client.create_collection(
            collection_name=collection_name,
            vectors_config=VectorParams(size=1024, distance=Distance.COSINE),
        )
    return collection_name


def index_chunks(book_id: str, chunks: list[str], metadatas: list[dict]):
    """Индексирует чанки в Qdrant с предварительной эмбеддизацией."""
    collection_name = get_or_create_collection(book_id)
    embeddings = _embed_documents(chunks)

    points = []
    for i, (chunk, embedding, metadata) in enumerate(zip(chunks, embeddings, metadatas)):
        point_id = str(uuid.uuid5(uuid.NAMESPACE_DNS, f"{book_id}_{i}"))
        points.append(
            PointStruct(
                id=point_id,
                vector=embedding,
                payload={"text": chunk, **metadata},
            )
        )

    client.upsert(collection_name=collection_name, points=points)


def delete_book_index(book_id: str):
    """Удаляет коллекцию книги (например, перед переиндексацией)."""
    collection_name = f"book_{book_id}"
    if client.collection_exists(collection_name=collection_name):
        client.delete_collection(collection_name=collection_name)


# ═══════════════════════════════════════════════════════════════
# Поиск
# ═══════════════════════════════════════════════════════════════

def search(book_id: str, query: str, top_k: int = 3) -> list[dict]:
    """Простой векторный поиск по Qdrant (без графа/реранка)."""
    collection_name = f"book_{book_id}"
    if not client.collection_exists(collection_name=collection_name):
        return []

    info = client.get_collection(collection_name=collection_name)
    if info.points_count == 0:
        return []

    query_vec = _embed_query(query)

    hits = client.query_points(
        collection_name=collection_name,
        query=query_vec,
        limit=top_k,
    ).points

    if not hits:
        return []

    return [
        {"text": hit.payload.get("text", ""), "metadata": hit.payload}
        for hit in hits
    ]


def hybrid_search(book_id: str, query: str, top_k: int = 5) -> list[dict]:
    """
    Гибридный поиск:
    1. Широкий векторный поиск (top‑20 из Qdrant).
    2. Опционально добавляются результаты из графа (если он построен).
    3. Дедупликация по началу текста.
    4. Ранжирование кросс‑энкодером до top_k.
    """
    # ── 1. Векторный поиск ──
    qdrant_results = search(book_id, query, top_k=20)

    # ── 2. Граф (если есть) ──
    graph_results = search_graph(book_id, query, top_k=10, query_threshold=0.5)

    # ── 3. Объединение и дедупликация ──
    seen = set()
    combined = []

    # Приоритет – граф (если он вообще что‑то дал)
    for r in graph_results:
        key = r["text"][:50].strip()
        if key and key not in seen:
            seen.add(key)
            combined.append({
                "text": r["text"],
                "metadata": {"page": r["page"], "book_id": book_id},
            })

    # Добавляем Qdrant‑результаты, избегая дубликатов
    for r in qdrant_results:
        key = r["text"][:50].strip()
        if key and key not in seen:
            seen.add(key)
            combined.append(r)

    # Если после объединения ничего нет – возвращаем пустоту
    if not combined:
        return []

    # ── 4. Ранжирование кросс‑энкодером ──
    return rerank(query, combined, top_k)