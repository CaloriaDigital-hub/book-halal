# stores.py
import os
from qdrant_client import QdrantClient
from qdrant_client.models import Distance, VectorParams, PointStruct
from interfaces import BaseVectorStore


class QdrantVectorStore(BaseVectorStore):
    def __init__(self, url: str | None = None):
        self.client = QdrantClient(url=url or os.getenv("QDRANT_URL", "http://localhost:6333"), check_compatibility=False)

    def create_collection(self, collection_name: str, vector_size: int):
        if not self.client.collection_exists(collection_name=collection_name):
            self.client.create_collection(
                collection_name=collection_name,
                vectors_config=VectorParams(size=vector_size, distance=Distance.COSINE),
            )

    def upsert(self, collection_name: str, points: list[dict]):
        qdrant_points = [
            PointStruct(
                id=p["id"],
                vector=p["vector"],
                payload=p["payload"],
            )
            for p in points
        ]
        self.client.upsert(collection_name=collection_name, points=qdrant_points)

    def query(self, collection_name: str, query_vector: list[float], top_k: int) -> list[dict]:
        if not self.collection_exists(collection_name):
            return []
        info = self.client.get_collection(collection_name=collection_name)
        if info.points_count == 0:
            return []

        hits = self.client.query_points(
            collection_name=collection_name,
            query=query_vector,
            limit=top_k,
        ).points

        return [
            {
                "text": hit.payload.get("text", ""),
                "metadata": hit.payload,
                "score": hit.score,
            }
            for hit in hits
        ]

    def delete_collection(self, collection_name: str):
        if self.client.collection_exists(collection_name=collection_name):
            self.client.delete_collection(collection_name=collection_name)

    def collection_exists(self, collection_name: str) -> bool:
        return self.client.collection_exists(collection_name=collection_name)

    def get_pages_by_numbers(self, collection_name: str, book_id: str, page_numbers: list[int]) -> list[dict]:
        """Возвращает все чанки указанных страниц (если есть)."""
        if not self.collection_exists(collection_name):
            return []
        results = []
        for page in page_numbers:
            points, _ = self.client.scroll(
                collection_name=collection_name,
                scroll_filter={
                    "must": [
                        {"key": "book_id", "match": {"value": book_id}},
                        {"key": "page", "match": {"value": page}}
                    ]
                },
                limit=100,
            )
            for p in sorted(points, key=lambda p: p.payload.get("chunk", 0)):
                results.append({
                    "text": p.payload.get("text", ""),
                    "metadata": p.payload,
                    "score": 0.0,          # будет пересчитан реранкером
                })
        return results