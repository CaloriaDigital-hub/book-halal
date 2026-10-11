# retrievers.py
import time
import logging
import torch
from sentence_transformers import CrossEncoder
from interfaces import BaseEmbedder, BaseVectorStore, BaseRetriever
from graph_legacy import search_graph

logger = logging.getLogger(__name__)


class HybridRetriever(BaseRetriever):
    def __init__(self, embedder: BaseEmbedder, vector_store: BaseVectorStore,
                 reranker_model_name: str = "BAAI/bge-reranker-v2-m3"):
        self.embedder = embedder
        self.vector_store = vector_store
        logger.info(f"Loading reranker model {reranker_model_name}...")
        t_load = time.time()
        self.reranker = CrossEncoder(reranker_model_name, max_length=512, device="cuda",
                                     automodel_args={"torch_dtype": torch.float16})
        logger.info(f"Reranker loaded in {time.time() - t_load:.2f}s")

    def retrieve(self, book_id: str, query: str, top_k: int = 5) -> list[dict]:
        collection_name = f"book_{book_id}"

        t_check = time.time()
        if not self.vector_store.collection_exists(collection_name):
            logger.warning(f"Collection {collection_name} does not exist")
            return []
        logger.info(f"[TIMING] collection check: {time.time() - t_check:.4f}s")

        # 1. Векторный поиск
        t_vec_start = time.time()
        query_vec = self.embedder.embed_query(query)
        logger.info(f"[TIMING] embed query: {time.time() - t_vec_start:.3f}s")

        t_qdrant = time.time()
        qdrant_results = self.vector_store.query(collection_name, query_vec, top_k=10)
        logger.info(f"[TIMING] qdrant query: {time.time() - t_qdrant:.3f}s, got {len(qdrant_results)} results")

        # 2. Графовый поиск (пока оставляем, но можно отключить при необходимости)
        t_graph_start = time.time()
        graph_results = search_graph(book_id, query, top_k=10, query_threshold=0.5)
        logger.info(f"[TIMING] graph search: {time.time() - t_graph_start:.3f}s, got {len(graph_results)} results")

        # 3. Объединение и дедупликация
        t_merge = time.time()
        seen = set()
        combined = []

        for r in graph_results:
            key = r["text"][:50].strip()
            if key and key not in seen:
                seen.add(key)
                combined.append({
                    "text": r["text"],
                    "metadata": {"page": r["page"], "book_id": book_id},
                    "score": r.get("score", 0.0),
                })

        for r in qdrant_results:
            key = r["text"][:50].strip()
            if key and key not in seen:
                seen.add(key)
                combined.append(r)
        logger.info(f"[TIMING] merge results: {time.time() - t_merge:.4f}s, combined {len(combined)} unique")

        # 3.5. Добавление соседних страниц (предыдущая и следующая)
        t_neighbor = time.time()
        pages_in_combined = set()
        for r in combined:
            page = r["metadata"].get("page")
            if page is not None:
                pages_in_combined.add(int(page))

        neighbor_pages = set()
        for page in pages_in_combined:
            if page > 1:
                neighbor_pages.add(page - 1)
            neighbor_pages.add(page + 1)  # страницы за пределами книги просто не найдутся
        neighbor_pages -= pages_in_combined

        if neighbor_pages:
            neighbors = self.vector_store.get_pages_by_numbers(
                collection_name, book_id, list(neighbor_pages)
            )
            for n in neighbors:
                key = n["text"][:50].strip()
                if key and key not in seen:
                    seen.add(key)
                    combined.append(n)
            logger.info(f"[TIMING] neighbor pages: {time.time() - t_neighbor:.3f}s, added {len(neighbors)} pages")
        else:
            logger.info(f"[TIMING] neighbor pages: {time.time() - t_neighbor:.3f}s, no new pages")

        if not combined:
            return []

        # 4. Реранкинг   
        t_rerank = time.time()
        texts = [doc["text"] for doc in combined]
        scores = self.reranker.predict([[query, t] for t in texts])
        logger.info(f"[TIMING] rerank predict ({len(texts)} pairs): {time.time() - t_rerank:.3f}s")

        for doc, score in zip(combined, scores):
            doc["score"] = float(score)

        combined.sort(key=lambda x: x["score"], reverse=True)
        final = combined[:top_k]

        logger.info(f"[TIMING] total retrieve: {time.time() - t_vec_start:.3f}s")
        return final