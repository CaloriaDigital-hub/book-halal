# interfaces.py
from abc import ABC, abstractmethod

class BaseEmbedder(ABC):
    @abstractmethod
    def embed_documents(self, texts: list[str]) -> list[list[float]]:
        """Преобразует список текстов в список векторов."""
        ...

    @abstractmethod
    def embed_query(self, query: str) -> list[float]:
        """Преобразует поисковый запрос в вектор."""
        ...

class BaseVectorStore(ABC):
    @abstractmethod
    def create_collection(self, collection_name: str, vector_size: int):
        """Создаёт коллекцию (индекс) для книги."""
        ...

    @abstractmethod
    def upsert(self, collection_name: str, points: list[dict]):
        """Вставляет или обновляет векторы с метаданными."""
        ...

    @abstractmethod
    def query(self, collection_name: str, query_vector: list[float], top_k: int) -> list[dict]:
        """Ищет top_k ближайших векторов. Возвращает список dict с ключами 'text', 'metadata', 'score'."""
        ...

    @abstractmethod
    def delete_collection(self, collection_name: str):
        """Удаляет коллекцию."""
        ...

    @abstractmethod
    def collection_exists(self, collection_name: str) -> bool:
        """Проверяет существование коллекции."""
        ...

class BaseRetriever(ABC):
    @abstractmethod
    def retrieve(self, book_id: str, query: str, top_k: int) -> list[dict]:
        """
        Выполняет поиск релевантных фрагментов по книге.
        Возвращает список dict с ключами 'text', 'metadata', 'score'.
        """
        ...

class BaseGenerator(ABC):
    @abstractmethod
    def generate(self, context: str, question: str) -> str:
        """
        Генерирует ответ строго на основе контекста.
        Принимает контекст (строка с фрагментами) и вопрос, возвращает ответ.
        """
        ...