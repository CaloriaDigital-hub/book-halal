import uuid
import os
from qdrant_client import QdrantClient
from qdrant_client.models import Distance, VectorParams, PointStruct
from sentence_transformers import SentenceTransformer

from graph import load_graph, search_graph
from sklearn.metrics.pairwise import cosine_similarity as cos_sim

model = SentenceTransformer("sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2")

# Подключаемся к Qdrant (по умолчанию http://localhost:6333)
client = QdrantClient(url=os.getenv("QDRANT_URL", "http://localhost:6333"))


def get_or_create_collection(book_id: str) -> str:
    collection_name = f"book_{book_id}"
    
    if not client.collection_exists(collection_name=collection_name):
        client.create_collection(
            collection_name=collection_name,
            vectors_config=VectorParams(size=384, distance=Distance.COSINE),
        )
    return collection_name


def index_chunks(book_id: str, chunks: list[str], metadatas: list[dict]):
    collection_name = get_or_create_collection(book_id)
    embeddings = model.encode(chunks).tolist()

    points = []
    for i, (chunk, embedding, metadata) in enumerate(zip(chunks, embeddings, metadatas)):
        # Генерируем детерминированный UUID на основе ID книги и индекса чанка
        point_id = str(uuid.uuid5(uuid.NAMESPACE_DNS, f"{book_id}_{i}"))
        
        points.append(
            PointStruct(
                id=point_id,
                vector=embedding,
                payload={"text": chunk, **metadata},
            )
        )

    # В Qdrant upsert тоже безопасен для многократного вызова
    client.upsert(
        collection_name=collection_name,
        points=points
    )


def search(book_id: str, query: str, top_k: int = 3) -> list[dict]:
    collection_name = f"book_{book_id}"

    # Если коллекции нет, значит ничего не проиндексировано
    if not client.collection_exists(collection_name=collection_name):
        return []

    # get collection record count to see if it's empty
    collection_info = client.get_collection(collection_name=collection_name)
    if collection_info.points_count == 0:
        return []

    query_embedding = model.encode([query])[0].tolist()
    
    search_results = client.query_points(
        collection_name=collection_name,
        query=query_embedding,
        limit=top_k,
    ).points

    if not search_results:
        return []

    return [
        {"text": hit.payload.get("text", ""), "metadata": hit.payload}
        for hit in search_results
    ]


def delete_book_index(book_id: str):
    """Delete all indexed chunks for a book (e.g., on re-upload)."""
    collection_name = f"book_{book_id}"
    if client.collection_exists(collection_name=collection_name):
        client.delete_collection(collection_name=collection_name)

def hybrid_search(book_id: str, query: str, top_k: int = 5) -> list[dict]:
    """
    Гибридный поиск: Qdrant кандидаты + граф переранжирование.
    """
    # Шаг 1 — Qdrant топ-10
    qdrant_results = search(book_id, query, top_k=10)
    if not qdrant_results:
        return []

    # Шаг 2 — граф соседи
    graph_results = search_graph(book_id, query, top_k=10, query_threshold=0.5)
    
    # Если графа нет — возвращаем просто Qdrant результаты
    if not graph_results:
        return qdrant_results[:top_k]

    # Шаг 3 — объединяем и дедуплицируем по тексту
    seen = set()
    combined = []

    # Сначала граф (он точнее по смыслу)
    for r in graph_results:
        key = r["text"][:50]
        if key not in seen:
            seen.add(key)
            combined.append({
                "text": r["text"],
                "metadata": {"page": r["page"], "book_id": book_id},
                "score": r["score"]
            })

    # Потом Qdrant добавляем то чего нет в графе
    for r in qdrant_results:
        key = r["text"][:50]
        if key not in seen:
            seen.add(key)
            combined.append(r)

    return combined[:top_k]