package repository

import (
	"context"
	"fmt"
	"golibrary/internal/entities"
	"time"

	"github.com/jmoiron/sqlx"
)

type BookLoanRepository interface {
	Create(ctx context.Context, bookID, userID int) (entities.BookLoan, error)
	FindActiveByBook(ctx context.Context, bookID int) (*entities.BookLoan, error)
	FindByUser(ctx context.Context, userID int) ([]entities.BookLoan, error)
	MarkReturned(ctx context.Context, loanID int) error
}

type bookLoanRepo struct {
	db *sqlx.DB
}

func NewBookLoanRepository(db *sqlx.DB) BookLoanRepository {
	return &bookLoanRepo{db: db}
}

func (r *bookLoanRepo) Create(ctx context.Context, bookID, userID int) (entities.BookLoan, error) {
	var id int
	query := `
        INSERT INTO book_loans (book_id, user_id)
        VALUES ($1, $2)
        RETURNING id
    `
	err := r.db.QueryRowContext(ctx, query, bookID, userID).Scan(&id)
	if err != nil {
		return entities.BookLoan{}, fmt.Errorf("failed to create loan: %w", err)
	}

	return r.getByID(ctx, id)
}

func (r *bookLoanRepo) getByID(ctx context.Context, id int) (entities.BookLoan, error) {
	var loan entities.BookLoan
	query := `
        SELECT id, book_id, user_id, borrowed_at, returned_at
        FROM book_loans
        WHERE id = $1
    `
	err := r.db.GetContext(ctx, &loan, query, id)
	if err != nil {
		return entities.BookLoan{}, fmt.Errorf("failed to get loan by id: %w", err)
	}

	return loan, nil
}

func (r *bookLoanRepo) FindActiveByBook(ctx context.Context, bookID int) (*entities.BookLoan, error) {
	var loan entities.BookLoan
	query := `
        SELECT id, book_id, user_id, borrowed_at, returned_at
        FROM book_loans
        WHERE book_id = $1 AND returned_at IS NULL
        LIMIT 1
    `
	err := r.db.GetContext(ctx, &loan, query, bookID)
	if err != nil {
		return nil, nil
	}

	return &loan, nil
}

func (r *bookLoanRepo) FindByUser(ctx context.Context, userID int) ([]entities.BookLoan, error) {
	var loans []entities.BookLoan
	query := `
        SELECT id, book_id, user_id, borrowed_at, returned_at
        FROM book_loans
        WHERE user_id = $1
        ORDER BY borrowed_at DESC
    `
	err := r.db.SelectContext(ctx, &loans, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to select loans by user: %w", err)
	}

	return loans, nil
}

func (r *bookLoanRepo) MarkReturned(ctx context.Context, loanID int) error {
	query := `
        UPDATE book_loans
        SET returned_at = $1
        WHERE id = $2 AND returned_at IS NULL
    `
	now := time.Now()
	res, err := r.db.ExecContext(ctx, query, now, loanID)
	if err != nil {
		return fmt.Errorf("failed to mark loan as returned: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no active loan found to mark as returned")
	}

	return nil
}
