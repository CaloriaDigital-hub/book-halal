export type RagHealth = { status: string };

export type IndexChunk = { page: number | null; chunk?: number | null; text: string | null };

export type IndexSnapshot = {
  book_id: string;
  total: number;
  chunks: IndexChunk[];
};

/** image_url должен быть абсолютным — RAG сам скачивает картинку и делает OCR. */
export type IndexPageInput = { page_number: number; image_url: string };

/** indexed/skipped — страницы, chunks — сколько чанков получилось после нарезки. */
export type IndexResult = { indexed: number; skipped: number; chunks?: number };

export type AskSource = { page: number; text: string };

export type AskResult = { answer: string; sources: AskSource[] };

export type PageChunk = { found: false } | { found: true; page: number; chunks?: number; text: string };
