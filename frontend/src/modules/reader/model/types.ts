export type ReaderPage = { page_number: number; image_url: string };

export type Bookmark = { id: string; page_number: number; created_at: string };

export type AskSource = { page: number; text: string };

export type AskAnswer = { answer: string; sources: AskSource[] | null };
