package db

import (
	"context"
	"golibrary/internal/entities"
	"golibrary/internal/repository"
	"log"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jmoiron/sqlx"
)

func SeedData(dbConn *sqlx.DB) error {
	authorRepo := repository.NewAuthorRepository(dbConn)
	userRepo := repository.NewUserRepository(dbConn)
	bookRepo := repository.NewBookRepository(dbConn)
	loanRepo := repository.NewBookLoanRepository(dbConn)

	countAuthors := getCount(dbConn, "authors")
	if countAuthors == 0 {
		log.Println("Seeding authors...")
		for i := 0; i < 10; i++ {
			author := entities.Author{Name: gofakeit.Name()}
			if _, err := authorRepo.Create(context.Background(), author); err != nil {
				return err
			}
		}
	}

	countUsers := getCount(dbConn, "users")
	if countUsers == 0 {
		log.Println("Seeding users...")
		for i := 0; i < 50; i++ {
			user := entities.User{Name: gofakeit.Name(), Email: gofakeit.Email()}
			if _, err := userRepo.Create(context.Background(), user); err != nil {
				return err
			}
		}
	}

	countBooks := getCount(dbConn, "books")
	if countBooks == 0 {
		log.Println("Seeding books...")
		authors, err := authorRepo.GetAll(context.Background())
		if err != nil {
			return err
		}
		for i := 0; i < 100; i++ {
			author := authors[gofakeit.Number(0, len(authors)-1)]
			book := entities.Book{
				Title:    gofakeit.BookTitle(),
				AuthorID: author.ID,
			}
			if _, err := bookRepo.Create(context.Background(), book); err != nil {
				return err
			}
		}
	}

	countLoans := getCount(dbConn, "book_loans")
	if countLoans == 0 {
		log.Println("Seeding book loans...")
		books, err := bookRepo.GetAll(context.Background())
		if err != nil {
			return err
		}
		users, err := userRepo.GetAll(context.Background())
		if err != nil {
			return err
		}
		for i := 0; i < 5; i++ {
			book := books[gofakeit.Number(0, len(books)-1)]
			user := users[gofakeit.Number(0, len(users)-1)]

			activeLoan, err := loanRepo.FindActiveByBook(context.Background(), book.ID)
			if err != nil {
				return err
			}
			if activeLoan == nil {
				if _, err := loanRepo.Create(context.Background(), book.ID, user.ID); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func getCount(db *sqlx.DB, tableName string) int {
	var count int
	query := "SELECT COUNT(*) FROM " + tableName
	if err := db.GetContext(context.Background(), &count, query); err != nil {
		log.Fatalf("failed to get count from %s: %v", tableName, err)
	}
	return count
}
