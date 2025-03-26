package service

import (
	"context"
	"errors"
	"golibrary/internal/entities"
	"golibrary/internal/repository"
)

type BookLoanService interface {
	BorrowBook(ctx context.Context, bookID, userID int) (entities.BookLoan, error)
	ReturnBook(ctx context.Context, bookID, userID int) error
	GetLoansByUser(ctx context.Context, userID int) ([]entities.BookLoan, error)
}

type bookLoanService struct {
	loansRepo repository.BookLoanRepository
	bookRepo  repository.BookRepository
}

func NewBookLoanService(loansRepo repository.BookLoanRepository, bookRepo repository.BookRepository) BookLoanService {
	return &bookLoanService{loansRepo: loansRepo, bookRepo: bookRepo}
}

func (s *bookLoanService) BorrowBook(ctx context.Context, bookID, userID int) (entities.BookLoan, error) {

	loan, err := s.loansRepo.FindActiveByBook(ctx, bookID)
	if err != nil {
		return entities.BookLoan{}, err
	}
	if loan != nil {
		return entities.BookLoan{}, errors.New("the book has already been issued")
	}

	return s.loansRepo.Create(ctx, bookID, userID)
}

func (s *bookLoanService) ReturnBook(ctx context.Context, bookID, userID int) error {
	loans, err := s.loansRepo.FindByUser(ctx, userID)
	if err != nil {
		return err
	}

	var targetLoanID int

	found := false

	for _, loan := range loans {
		if loan.BookID == bookID && loan.ReturnedAt == nil {
			targetLoanID = loan.ID
			found = true
			break
		}
	}
	if !found {
		return errors.New("user does not have an active loan for this book")
	}

	return s.loansRepo.MarkReturned(ctx, targetLoanID)
}

func (s *bookLoanService) GetLoansByUser(ctx context.Context, userID int) ([]entities.BookLoan, error) {
	return s.loansRepo.FindByUser(ctx, userID)
}
