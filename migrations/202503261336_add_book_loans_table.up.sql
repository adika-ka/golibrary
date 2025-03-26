CREATE TABLE IF NOT EXISTS book_loans (
    id SERIAL PRIMARY KEY,
    book_id INT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    borrowed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    returned_at TIMESTAMP NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_book_loans_book_active 
    ON book_loans(book_id)
    WHERE returned_at IS NULL;

ALTER TABLE books DROP COLUMN IF EXISTS rented_by;