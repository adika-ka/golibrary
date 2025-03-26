CREATE TABLE IF NOT EXISTS books (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    author_id INT NOT NULL REFERENCES authors(id) ON DELETE CASCADE,
    rented_by INT REFERENCES users(id)
);