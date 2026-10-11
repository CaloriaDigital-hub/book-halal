-- Покупки книг: чек загружает покупатель, подтверждает админ (в Telegram или в админке)
CREATE TABLE purchases (
    id                  UUID PRIMARY KEY,
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    book_id             UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    amount              INT NOT NULL,                      -- цена книги на момент покупки, в минорных единицах (как books.price)
    status              VARCHAR(20) NOT NULL DEFAULT 'pending',
    contact_type        VARCHAR(20) NOT NULL,              -- phone | whatsapp | telegram
    contact_value       VARCHAR(100) NOT NULL,
    receipt_path        TEXT NOT NULL,                     -- путь в приватном хранилище, не в static
    receipt_mime        VARCHAR(100) NOT NULL,
    telegram_message_id BIGINT,
    reviewed_by         VARCHAR(255),
    reviewed_at         TIMESTAMP WITH TIME ZONE,
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Одна активная (ожидает/подтверждена) покупка на пару (user, book); после отказа можно отправить чек заново
CREATE UNIQUE INDEX unique_active_purchase ON purchases(user_id, book_id) WHERE status IN ('pending', 'approved');
CREATE INDEX idx_purchases_status_created ON purchases(status, created_at DESC);

-- Реквизиты для оплаты: одна строка на весь сервис
CREATE TABLE payment_settings (
    id         SMALLINT PRIMARY KEY DEFAULT 1 CHECK (id = 1),
    details    TEXT NOT NULL DEFAULT '',
    qr_url     TEXT NOT NULL DEFAULT '',
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

INSERT INTO payment_settings (id) VALUES (1);
