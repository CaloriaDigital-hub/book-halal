-- Удаляем таблицы в обратном порядке (сначала зависимую, потом главную)
DROP TABLE IF EXISTS book_pages;
DROP TABLE IF EXISTS books;