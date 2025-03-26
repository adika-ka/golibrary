package service

import (
	"context"
	"fmt"
	"golibrary/internal/entities"
	"golibrary/internal/repository"
)

type AuthorService interface {
	CreateAuthor(ctx context.Context, name string) (entities.Author, error)
	ListAuthors(ctx context.Context) ([]entities.Author, error)
	TopAuthors(ctx context.Context, limit int) ([]entities.TopAuthor, error)
	GetByID(ctx context.Context, id int) (entities.Author, error)
}

type authorService struct {
	repo     repository.AuthorRepository
	bookRepo repository.BookRepository
}

func NewAuthorService(repo repository.AuthorRepository, bookRepo repository.BookRepository) AuthorService {
	return &authorService{repo: repo, bookRepo: bookRepo}
}

func (s *authorService) CreateAuthor(ctx context.Context, name string) (entities.Author, error) {
	if name == "" {
		return entities.Author{}, fmt.Errorf("author name cannot be empty")
	}

	return s.repo.Create(ctx, entities.Author{Name: name})
}

func (s *authorService) ListAuthors(ctx context.Context) ([]entities.Author, error) {
	authors, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for i, author := range authors {
		books, err := s.bookRepo.FindByAuthor(ctx, author.ID)
		if err != nil {
			return nil, err
		}
		authors[i].Books = books
	}

	return authors, nil
}

func (s *authorService) TopAuthors(ctx context.Context, limit int) ([]entities.TopAuthor, error) {
	return s.repo.GetTopAuthors(ctx, limit)
}

func (s *authorService) GetByID(ctx context.Context, id int) (entities.Author, error) {
	author, err := s.repo.Get(ctx, id)
	if err != nil {
		return entities.Author{}, err
	}

	books, err := s.bookRepo.FindByAuthor(ctx, author.ID)
	if err != nil {
		return entities.Author{}, err
	}

	author.Books = books
	return author, nil
}
