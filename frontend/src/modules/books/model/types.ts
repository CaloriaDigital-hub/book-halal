export type BookCatalogItem = {
  id: string;
  title: string;
  author: string;
  description: string;
  price: number;
  cover_url: string;
};

export type BookStatus = "processing" | "ready" | "error";

export type BookDetails = BookCatalogItem & {
  total_pages: number;
  status: BookStatus | string;
};

export type BookPage = {
  page_number: number;
  image_url: string;
};

export type BookPages = {
  total_pages: number;
  pages: BookPage[] | null;
};

export type UploadBookInput = {
  title: string;
  author: string;
  description: string;
  price: number;
  document: File;
};

export type UploadBookResult = {
  message: string;
  book_id: string;
};
