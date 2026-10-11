# generators.py
import os
from openai import OpenAI
from interfaces import BaseGenerator

SYSTEM_PROMPT = """Ты — ассистент, который объясняет материал исламской книги, опираясь ТОЛЬКО на предоставленный контекст.
Твоя задача — дать полезный, понятный ответ, как если бы ты пересказывал прочитанное своими словами, но без искажения смысла.

ПРАВИЛА:
1. Используй ТОЛЬКО информацию из контекста. Никаких внешних знаний.
2. Если в контексте нет точного ответа, напиши: "В данной книге ответ не найден." и остановись.
3. Отвечай на том же языке, что и вопрос.
4. Перефразируй и комбинируй фрагменты контекста так, чтобы ответ звучал связно и естественно, но сохраняй точность фактов.
5. В конце ответа обязательно укажи источник (или источники) в формате [Источник: страница X]. Если источников несколько, перечисли их.
6. Не добавляй ничего от себя — ни примеров, ни объяснений, ни сравнений с другими книгами или мнениями.
7. Если в контексте приведены аяты или хадисы, их нужно процитировать дословно.
8. Отвечай развёрнуто, но по существу вопроса, избегай воды."""

class DeepSeekGenerator(BaseGenerator):
    def __init__(self, api_key: str | None = None, base_url: str = "https://api.deepseek.com",
                 model: str = "deepseek-chat", temperature: float = 0.1):
        self.client = OpenAI(
            api_key=api_key or os.getenv("LLM_API_KEYS"),
            base_url=base_url,
        )
        self.model = model
        self.temperature = temperature

    def generate(self, context: str, question: str) -> str:
        """
        Принимает контекст (строка с уже размеченными фрагментами) и вопрос.
        Возвращает ответ модели.
        """
        response = self.client.chat.completions.create(
            model=self.model,
            temperature=self.temperature,
            messages=[
                {"role": "system", "content": SYSTEM_PROMPT},
                {"role": "user", "content": f"Контекст:\n---\n{context}\n---\n\nВопрос: {question}"},
            ],
        )
        return response.choices[0].message.content