-- Создаем таблицу для книг
CREATE TABLE books (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    description TEXT,                               
    price INT NOT NULL DEFAULT 0,                   
    cover_url VARCHAR(255),                     
    total_pages INT NOT NULL DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'processing',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);


CREATE TABLE book_pages (
    id UUID PRIMARY KEY,
    book_id UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    page_number INT NOT NULL,
    image_url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
   
    CONSTRAINT unique_book_page UNIQUE(book_id, page_number)
);ц


CREATE INDEX idx_book_pages_book_id ON book_pages(book_id);