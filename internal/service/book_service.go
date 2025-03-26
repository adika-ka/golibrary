package service

import (
	"context"
	"fmt"
	"golibrary/internal/entities"
	"golibrary/internal/repository"
)

type BookService interface {
	CreateBook(ctx context.Context, book entities.Book) (entities.Book, error)
	FindBookByID(ctx context.Context, id int) (entities.Book, error)
	ListBooks(ctx context.Context) ([]entities.Book, error)
}

type bookService struct {
	repo repository.BookRepository
}

func NewBookService(repo repository.BookRepository) BookService {
	return &bookService{repo: repo}
}

func (s *bookService) CreateBook(ctx context.Context, book entities.Book) (entities.Book, error) {
	if book.Title == "" {
		return entities.Book{}, fmt.Errorf("title cannot be empty")
	}

	if book.AuthorID <= 0 {
		return entities.Book{}, fmt.Errorf("author id cannot be empty")
	}

	return s.repo.Create(ctx, book)
}

func (s *bookService) FindBookByID(ctx context.Context, id int) (entities.Book, error) {
	if id <= 0 {
		return entities.Book{}, fmt.Errorf("incorrect id %d", id)
	}

	return s.repo.GetByID(ctx, id)
}

func (s *bookService) ListBooks(ctx context.Context) ([]entities.Book, error) {
	return s.repo.GetAll(ctx)
}
