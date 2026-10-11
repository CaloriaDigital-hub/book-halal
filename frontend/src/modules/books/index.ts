// Публичный API модуля books. Всё, что не экспортировано здесь, — внутреннее.
export { getAllPages } from "./api/books-api";
export type { BookPage } from "./model/types";
export { BookSummaryWidget } from "./ui/BookSummaryWidget";
export { BooksCatalogWidget } from "./ui/BooksCatalogWidget";
export { CatalogStatsWidget } from "./ui/CatalogStatsWidget";
export { UploadBookWidget } from "./ui/UploadBookWidget";
