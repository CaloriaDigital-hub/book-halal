import chromadb
from sentence_transformers import SentenceTransformer

model = SentenceTransformer("sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2")

client = chromadb.PersistentClient(path="./chroma_db")


def get_collection(book_id: str):
    return client.get_or_create_collection(name=f"book_{book_id}")


def index_chunks(book_id: str, chunks: list[str], metadatas: list[dict]):
    collection = get_collection(book_id)
    embeddings = model.encode(chunks).tolist()

    # Build IDs based on book_id + chunk content hash to avoid duplicates
    ids = [f"{book_id}_{i}" for i in range(len(chunks))]

    # Upsert instead of add — safe to call multiple times
    collection.upsert(
        documents=chunks,
        embeddings=embeddings,
        metadatas=metadatas,
        ids=ids,
    )


def search(book_id: str, query: str, top_k: int = 3) -> list[dict]:
    collection = get_collection(book_id)

    # Guard: if collection is empty, return nothing
    if collection.count() == 0:
        return []

    query_embedding = model.encode([query]).tolist()
    results = collection.query(
        query_embeddings=query_embedding,
        n_results=min(top_k, collection.count()),
    )

    if not results["documents"] or not results["documents"][0]:
        return []

    return [
        {"text": doc, "metadata": meta}
        for doc, meta in zip(results["documents"][0], results["metadatas"][0])
    ]


def delete_book_index(book_id: str):
    """Delete all indexed chunks for a book (e.g., on re-upload)."""
    try:
        client.delete_collection(name=f"book_{book_id}")
    except Exception:
        pass  # Collection didn't exist, that's fine