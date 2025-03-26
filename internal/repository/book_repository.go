package repository

import (
	"context"
	"fmt"
	"golibrary/internal/entities"

	"github.com/jmoiron/sqlx"
)

type BookRepository interface {
	Create(ctx context.Context, book entities.Book) (entities.Book, error)
	GetByID(ctx context.Context, id int) (entities.Book, error)
	GetAll(ctx context.Context) ([]entities.Book, error)
	Update(ctx context.Context, id int, book entities.Book) (entities.Book, error)
	Delete(ctx context.Context, id int) error
	FindByAuthor(ctx context.Context, authorID int) ([]entities.Book, error)
}

type bookRepo struct {
	db *sqlx.DB
}

func NewBookRepository(db *sqlx.DB) BookRepository {
	return &bookRepo{db: db}
}

func (r *bookRepo) Create(ctx context.Context, book entities.Book) (entities.Book, error) {
	var id int

	query := `INSERT INTO books (title, author_id) VALUES ($1, $2) RETURNING id`

	err := r.db.QueryRowContext(ctx, query, book.Title, book.AuthorID).Scan(&id)
	if err != nil {
		return entities.Book{}, fmt.Errorf("failed to insert book: %w", err)
	}

	return r.GetByID(ctx, id)
}

func (r *bookRepo) GetByID(ctx context.Context, id int) (entities.Book, error) {
	var book entities.Book

	query := `SELECT id, title, author_id FROM books WHERE id = $1`

	err := r.db.GetContext(ctx, &book, query, id)
	if err != nil {
		return entities.Book{}, fmt.Errorf("failed to get book by id: %w", err)
	}

	return book, nil
}

func (r *bookRepo) GetAll(ctx context.Context) ([]entities.Book, error) {
	var books []entities.Book

	query := `SELECT id, title, author_id FROM books`

	err := r.db.SelectContext(ctx, &books, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all books: %w", err)
	}

	return books, nil
}

func (r *bookRepo) Update(ctx context.Context, id int, book entities.Book) (entities.Book, error) {
	query := `UPDATE books SET title = $1, author_id = $2 WHERE id = $3`

	_, err := r.db.ExecContext(ctx, query, book.Title, book.AuthorID, id)
	if err != nil {
		return entities.Book{}, fmt.Errorf("failed to update book: %w", err)
	}

	return r.GetByID(ctx, id)
}

func (r *bookRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM books WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete book: %w", err)
	}

	return nil
}

func (r *bookRepo) FindByAuthor(ctx context.Context, authorID int) ([]entities.Book, error) {
	var books []entities.Book

	query := `SELECT id, title, author_id FROM books WHERE author_id = $1`

	err := r.db.SelectContext(ctx, &books, query, authorID)
	if err != nil {
		return nil, fmt.Errorf("failed to find books by author: %w", err)
	}

	return books, nil
}
