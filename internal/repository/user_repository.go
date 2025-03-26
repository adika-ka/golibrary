package repository

import (
	"context"
	"database/sql"
	"fmt"
	"golibrary/internal/entities"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(ctx context.Context, user entities.User) (entities.User, error)
	GetByID(ctx context.Context, id int) (entities.User, error)
	GetAll(ctx context.Context) ([]entities.User, error)
	Update(ctx context.Context, id int, user entities.User) (entities.User, error)
	Delete(ctx context.Context, id int) error
	GetAllWithBooks(ctx context.Context) ([]entities.User, error)
	ExistsByID(ctx context.Context, id int) (bool, error)
}

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepo{db: db}
}

func (u *userRepo) Create(ctx context.Context, user entities.User) (entities.User, error) {
	query := `
		INSERT INTO users (name, email)
		VALUES ($1, $2)
		RETURNING id;
	`

	var newID int
	err := u.db.QueryRowContext(ctx, query, user.Name, user.Email).Scan(&newID)
	if err != nil {
		return user, fmt.Errorf("failed to insert user: %v", err)
	}

	user.ID = newID
	return user, nil
}

func (u *userRepo) GetByID(ctx context.Context, id int) (entities.User, error) {
	query := `
		SELECT id, name, email
		FROM users WHERE id = $1
	`

	var user entities.User

	err := u.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return user, fmt.Errorf("failed to find user by id %d: %w", id, err)
	}

	return user, nil
}

func (u *userRepo) GetAll(ctx context.Context) ([]entities.User, error) {
	query := `
		SELECT id, name, email
		FROM users
	`

	var users []entities.User
	err := u.db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, fmt.Errorf("failed to select users: %w", err)
	}

	return users, nil
}

func (u *userRepo) Update(ctx context.Context, id int, user entities.User) (entities.User, error) {
	query := `
		UPDATE users
		SET name = $1, email = $2
		WHERE id = $3
	`
	_, err := u.db.ExecContext(ctx, query,
		user.Name,
		user.Email,
		id,
	)

	if err != nil {
		return user, fmt.Errorf("failed to update user by id %d: %w", id, err)
	}

	updateUser, err := u.GetByID(ctx, id)
	if err != nil {
		return user, fmt.Errorf("failed to retrieve updated user %d: %w", id, err)
	}

	return updateUser, nil
}

func (u *userRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := u.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user %d: %w", id, err)
	}
	return nil
}

func (u *userRepo) GetAllWithBooks(ctx context.Context) ([]entities.User, error) {
	query := `
		SELECT 
			u.id AS user_id, 
			u.name AS user_name, 
			u.email AS user_email,
			b.id AS book_id, 
			b.title AS book_title, 
			b.author_id AS book_author_id,
			a.id AS author_id,
			a.name AS author_name
		FROM users u
		LEFT JOIN book_loans bl ON bl.user_id = u.id AND bl.returned_at IS NULL
		LEFT JOIN books b ON b.id = bl.book_id
		LEFT JOIN authors a ON a.id = b.author_id
		ORDER BY u.id
	`

	var rows []struct {
		UserID       int            `db:"user_id"`
		UserName     string         `db:"user_name"`
		UserEmail    string         `db:"user_email"`
		BookID       sql.NullInt64  `db:"book_id"`
		BookTitle    sql.NullString `db:"book_title"`
		BookAuthorID sql.NullInt64  `db:"book_author_id"`
		AuthorID     sql.NullInt64  `db:"author_id"`
		AuthorName   sql.NullString `db:"author_name"`
	}

	err := u.db.SelectContext(ctx, &rows, query)
	if err != nil {
		return nil, fmt.Errorf("failed to select users with books: %w", err)
	}

	userMap := make(map[int]*entities.User)
	for _, row := range rows {
		user, exists := userMap[row.UserID]
		if !exists {
			user = &entities.User{
				ID:          row.UserID,
				Name:        row.UserName,
				Email:       row.UserEmail,
				RentedBooks: []entities.Book{},
			}
			userMap[row.UserID] = user
		}

		if row.BookID.Valid {
			book := entities.Book{
				ID:       int(row.BookID.Int64),
				Title:    row.BookTitle.String,
				AuthorID: int(row.BookAuthorID.Int64),
			}

			if row.AuthorID.Valid {
				book.Author = &entities.Author{
					ID:   int(row.AuthorID.Int64),
					Name: row.AuthorName.String,
				}
			}
			user.RentedBooks = append(user.RentedBooks, book)
		}
	}

	users := make([]entities.User, 0, len(userMap))
	for _, user := range userMap {
		users = append(users, *user)
	}

	return users, nil
}

func (u *userRepo) ExistsByID(ctx context.Context, id int) (bool, error) {
	query := `SELECT 1 FROM users WHERE id = $1 LIMIT 1`

	var dummy int
	err := u.db.QueryRowContext(ctx, query, id).Scan(&dummy)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
