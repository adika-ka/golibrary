package facade

import (
	"context"
	"golibrary/internal/entities"
	"golibrary/internal/service"
)

type LibraryFacade struct {
	service *service.LibraryService
}

func NewLibraryFacade(service *service.LibraryService) *LibraryFacade {
	return &LibraryFacade{service: service}
}

func (f *LibraryFacade) RentBook(userID, bookID int) error {
	return f.service.RentBook(userID, bookID)
}

func (f *LibraryFacade) ReturnBook(userID, bookID int) error {
	return f.service.ReturnBook(userID, bookID)
}

func (f *LibraryFacade) GetAllBooks(ctx context.Context) ([]entities.Book, error) {
	return f.service.GetAllBooks(ctx)
}

func (f *LibraryFacade) CreateBook(ctx context.Context, book entities.Book) (entities.Book, error) {
	return f.service.CreateBook(ctx, book)
}

func (f *LibraryFacade) GetAllAuthors(ctx context.Context) ([]entities.Author, error) {
	return f.service.GetAllAuthors(ctx)
}

func (f *LibraryFacade) CreateAuthor(ctx context.Context, name string) (entities.Author, error) {
	return f.service.CreateAuthor(ctx, name)
}

func (f *LibraryFacade) GetTopAuthors(ctx context.Context, limit int) ([]entities.TopAuthor, error) {
	return f.service.GetTopAuthors(ctx, limit)
}

func (f *LibraryFacade) GetAllUsers(ctx context.Context) ([]entities.User, error) {
	return f.service.GetAllUsers(ctx)
}

func (f *LibraryFacade) GetBooksBorrowedByUser(ctx context.Context, userID int) ([]entities.Book, error) {
	return f.service.GetBooksBorrowedByUser(ctx, userID)
}
