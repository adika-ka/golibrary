package repository

import (
	"context"
	"fmt"
	"golibrary/internal/entities"

	"github.com/jmoiron/sqlx"
)

type AuthorRepository interface {
	Create(ctx context.Context, author entities.Author) (entities.Author, error)
	GetAll(ctx context.Context) ([]entities.Author, error)
	GetTopAuthors(ctx context.Context, limit int) ([]entities.TopAuthor, error)
	Get(ctx context.Context, id int) (entities.Author, error)
}

type authorRepo struct {
	db *sqlx.DB
}

func NewAuthorRepository(db *sqlx.DB) AuthorRepository {
	return &authorRepo{db: db}
}

func (a *authorRepo) Create(ctx context.Context, author entities.Author) (entities.Author, error) {
	query := `
		INSERT INTO authors (name)
		VALUES ($1)
		RETURNING id;
	`

	var newID int
	err := a.db.QueryRowContext(ctx, query, author.Name).Scan(&newID)
	if err != nil {
		return author, fmt.Errorf("failed to insert author: %w", err)
	}

	author.ID = newID
	return author, nil
}

func (a *authorRepo) GetAll(ctx context.Context) ([]entities.Author, error) {
	query := `
		SELECT id, name
		FROM authors
	`

	var authors []entities.Author
	err := a.db.SelectContext(ctx, &authors, query)
	if err != nil {
		return nil, fmt.Errorf("failed to select authors: %w", err)
	}

	return authors, nil
}

func (a *authorRepo) GetTopAuthors(ctx context.Context, limit int) ([]entities.TopAuthor, error) {
	query := `
		SELECT authors.id, authors.name, COUNT(bl.id) AS rented_books_count
		FROM authors
		JOIN books b ON b.author_id = authors.id
		JOIN book_loans bl ON bl.book_id = b.id
		WHERE bl.returned_at IS NULL
		GROUP BY authors.id
		ORDER BY rented_books_count DESC
		LIMIT $1
	`

	var topAuthors []entities.TopAuthor
	err := a.db.SelectContext(ctx, &topAuthors, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to select top authors: %w", err)
	}

	return topAuthors, nil
}

func (a *authorRepo) Get(ctx context.Context, id int) (entities.Author, error) {
	var author entities.Author

	query := `SELECT id, name FROM authors WHERE id = $1`

	err := a.db.GetContext(ctx, &author, query, id)
	if err != nil {
		return author, fmt.Errorf("failed to get author by id: %w", err)
	}

	return author, nil
}
