# embeddings.py
import logging
from sentence_transformers import SentenceTransformer
from interfaces import BaseEmbedder
import torch

logger = logging.getLogger(__name__)

class SentenceTransformerEmbedder(BaseEmbedder):
    def __init__(self, model_name: str = "intfloat/multilingual-e5-large"):
        self._model = SentenceTransformer(model_name, device="cuda")
        logger.info(f"Model device: {self._model.device}")
        logger.info(f"CUDA available: {torch.cuda.is_available()}")

    def embed_documents(self, texts: list[str]) -> list[list[float]]:
        """Эмбеддинги для индексации (с префиксом 'passage: ')."""
        return self._model.encode(
            [f"passage: {t}" for t in texts],
            normalize_embeddings=True,
        ).tolist()

    def embed_query(self, query: str) -> list[float]:
        """Эмбеддинг запроса (с префиксом 'query: ')."""
        return self._model.encode(
            f"query: {query}",
            normalize_embeddings=True,
        ).tolist()