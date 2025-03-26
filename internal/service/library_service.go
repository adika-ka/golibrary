package service

import (
	"context"
	"fmt"
	"golibrary/internal/entities"
)

type LibraryService struct {
	bookService     BookService
	userService     UserService
	authorService   AuthorService
	bookLoanService BookLoanService
}

func NewLibraryService(bookService BookService, userService UserService, authorService AuthorService, bookLoanService BookLoanService) *LibraryService {
	return &LibraryService{
		bookService:     bookService,
		userService:     userService,
		authorService:   authorService,
		bookLoanService: bookLoanService,
	}
}

func (s *LibraryService) RentBook(userID, bookID int) error {
	exists, err := s.userService.ExistsUserByID(context.Background(), userID)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("user with ID %d does not exist", userID)
	}

	_, err = s.bookService.FindBookByID(context.Background(), bookID)
	if err != nil {
		return err
	}

	_, err = s.bookLoanService.BorrowBook(context.Background(), bookID, userID)
	return err
}

func (s *LibraryService) ReturnBook(userID, bookID int) error {
	exists, err := s.userService.ExistsUserByID(context.Background(), userID)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("user with ID %d does not exist", userID)
	}

	_, err = s.bookService.FindBookByID(context.Background(), bookID)
	if err != nil {
		return err
	}

	return s.bookLoanService.ReturnBook(context.Background(), bookID, userID)
}

func (s *LibraryService) GetAllBooks(ctx context.Context) ([]entities.Book, error) {
	return s.bookService.ListBooks(ctx)
}

func (s *LibraryService) CreateBook(ctx context.Context, book entities.Book) (entities.Book, error) {
	authors, err := s.authorService.ListAuthors(ctx)
	if err != nil {
		return book, fmt.Errorf("failed to verify author: %w", err)
	}

	found := false
	for _, a := range authors {
		if a.ID == book.AuthorID {
			found = true
			break
		}
	}

	if !found {
		return book, fmt.Errorf("author with ID %d does not exist", book.AuthorID)
	}

	return s.bookService.CreateBook(ctx, book)
}

func (s *LibraryService) GetAllAuthors(ctx context.Context) ([]entities.Author, error) {
	return s.authorService.ListAuthors(ctx)
}

func (s *LibraryService) CreateAuthor(ctx context.Context, name string) (entities.Author, error) {
	return s.authorService.CreateAuthor(ctx, name)
}

func (s *LibraryService) GetTopAuthors(ctx context.Context, limit int) ([]entities.TopAuthor, error) {
	return s.authorService.TopAuthors(ctx, limit)
}

func (s *LibraryService) GetAllUsers(ctx context.Context) ([]entities.User, error) {
	return s.userService.ListUsers(ctx)
}

func (s *LibraryService) GetBooksBorrowedByUser(ctx context.Context, userID int) ([]entities.Book, error) {
	loans, err := s.bookLoanService.GetLoansByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	var books []entities.Book

	for _, loan := range loans {
		if loan.ReturnedAt == nil {
			book, err := s.bookService.FindBookByID(ctx, loan.BookID)

			if err == nil {
				books = append(books, book)
			}
		}
	}

	return books, nil
}
