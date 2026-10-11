# pipeline.py
import time
import uuid
import logging
from interfaces import BaseEmbedder, BaseVectorStore, BaseRetriever, BaseGenerator
from graph_legacy import build_graph, delete_graph

logger = logging.getLogger(__name__)


class BookRAG:
    def __init__(self, embedder: BaseEmbedder, vector_store: BaseVectorStore,
                 retriever: BaseRetriever, generator: BaseGenerator):
        self.embedder = embedder
        self.vector_store = vector_store
        self.retriever = retriever
        self.generator = generator

    def index(self, book_id: str, chunks: list[str], metadatas: list[dict]):
        """Индексация книги: векторы + граф."""
        collection_name = f"book_{book_id}"

        t_start = time.time()
        sample_vec = self.embedder.embed_documents([chunks[0]])[0]
        self.vector_store.create_collection(collection_name, len(sample_vec))
        logger.info(f"[TIMING] create_collection: {time.time() - t_start:.3f}s")

        t_emb = time.time()
        embeddings = self.embedder.embed_documents(chunks)
        logger.info(f"[TIMING] embed_documents ({len(chunks)} chunks): {time.time() - t_emb:.3f}s")

        points = []
        for i, (chunk, emb, meta) in enumerate(zip(chunks, embeddings, metadatas)):
            points.append({
                # Qdrant принимает только uint или UUID — детерминированный UUID из прежнего ключа
                "id": str(uuid.uuid5(uuid.NAMESPACE_URL, f"{book_id}_{i}")),
                "vector": emb,
                "payload": {"text": chunk, **meta},
            })

        t_upsert = time.time()
        self.vector_store.upsert(collection_name, points)
        logger.info(f"[TIMING] upsert {len(points)} points: {time.time() - t_upsert:.3f}s")

        pages = [m["page"] for m in metadatas]
        t_graph = time.time()
        build_graph(book_id, chunks, pages)
        logger.info(f"[TIMING] build_graph: {time.time() - t_graph:.3f}s")
        logger.info(f"[TIMING] total index: {time.time() - t_start:.3f}s")

    def ask(self, book_id: str, question: str) -> dict:
        """Запрос: поиск + генерация."""
        t_total = time.time()

        # Замер retrieval (весь)
        t_ret = time.time()
        results = self.retriever.retrieve(book_id, question, top_k=20)
        t_ret_end = time.time()
        logger.info(f"[TIMING] retrieval total: {t_ret_end - t_ret:.3f}s, found {len(results)} chunks")

        if not results:
            return {"answer": "В данной книге ответ не найден.", "sources": []}

        # Сборка контекста
        t_ctx = time.time()
        context_parts = []
        for i, r in enumerate(results, 1):
            page = r["metadata"]["page"]
            text = r["text"]
            context_parts.append(f"[Фрагмент {i}]\n{text}\n[Источник: страница {page}]\n")
        context = "\n".join(context_parts)
        logger.info(f"[TIMING] build context: {time.time() - t_ctx:.4f}s")

        # Генерация
        t_gen = time.time()
        answer = self.generator.generate(context, question)
        t_gen_end = time.time()
        logger.info(f"[TIMING] generation (LLM): {t_gen_end - t_gen:.3f}s")

        logger.info(f"[TIMING] total ask: {t_gen_end - t_total:.3f}s")

        if answer.startswith("В данной книге ответ не найден"):
            return {"answer": answer, "sources": []}

        return {
            "answer": answer,
            "sources": [
                {"page": r["metadata"]["page"], "text": r["text"]}
                for r in results[:5]
            ],
        }

    def delete_index(self, book_id: str):
        """Удаление индекса и графа."""
        collection_name = f"book_{book_id}"
        t_del = time.time()
        if self.vector_store.collection_exists(collection_name):
            self.vector_store.delete_collection(collection_name)
        delete_graph(book_id)
        logger.info(f"[TIMING] delete_index: {time.time() - t_del:.3f}s")