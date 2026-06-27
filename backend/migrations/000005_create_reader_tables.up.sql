-- Прогресс чтения: одна запись на пару (user, book)
CREATE TABLE reading_progress (
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    book_id     UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    page_number INT NOT NULL DEFAULT 1,
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, book_id)
);

-- Закладки: уникальна пара (user, book, page)
CREATE TABLE bookmarks (
    id          UUID PRIMARY KEY,
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    book_id     UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    page_number INT NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_bookmark UNIQUE(user_id, book_id, page_number)
);

CREATE INDEX idx_bookmarks_user_book ON bookmarks(user_id, book_id);
CREATE INDEX idx_reading_progress_user ON reading_progress(user_id);
