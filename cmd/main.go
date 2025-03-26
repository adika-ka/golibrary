package main

import (
	"context"
	_ "golibrary/docs"
	"golibrary/infrastructure"
	"golibrary/internal/controller"
	"golibrary/internal/db"
	"golibrary/internal/facade"
	"golibrary/internal/repository"
	"golibrary/internal/service"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-chi/chi/v5"
)

// @title Library API
// @version 1.0
// @description API для управления библиотекой.
// @host localhost:8080
// @BasePath /

func main() {
	dbConn, err := db.InitDBAndMigrate()
	if err != nil {
		log.Fatalf("failed to init DB and run migrations: %v", err)
	}
	defer dbConn.Close()

	err = db.SeedData(dbConn)
	if err != nil {
		log.Fatalf("failed to create initial data: %v", err)
	}

	r := chi.NewRouter()

	bookRepo := repository.NewBookRepository(dbConn)
	authorRepo := repository.NewAuthorRepository(dbConn)
	userRepo := repository.NewUserRepository(dbConn)
	loanRepo := repository.NewBookLoanRepository(dbConn)

	bookService := service.NewBookService(bookRepo)
	authorService := service.NewAuthorService(authorRepo, bookRepo)
	userService := service.NewUserService(userRepo)
	bookLoanService := service.NewBookLoanService(loanRepo, bookRepo)

	libraryService := service.NewLibraryService(bookService, userService, authorService, bookLoanService)
	libraryFacade := facade.NewLibraryFacade(libraryService)

	responder := infrastructure.NewJSONResponder()

	bookController := controller.NewBookController(libraryFacade, responder)
	authorController := controller.NewAuthorController(libraryFacade, responder)
	userController := controller.NewUserController(libraryFacade, responder)
	loanController := controller.NewLoanController(libraryFacade, responder)

	controller.RegisterBookRoutes(r, bookController)
	controller.RegisterAuthorRoutes(r, authorController)
	controller.RegisterUserRoutes(r, userController)
	controller.RegisterLoanRoutes(r, loanController)

	r.Post("/loans/borrow", loanController.BorrowBookEndpoint)
	r.Post("/loans/return", loanController.ReturnBookEndpoint)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatalf("Error creating listener: %v", err)
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("server is starting at :8080")
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()
	<-stopChan
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forsed to shutdown: %v", err)
	}
	log.Println("Server stopped gracefully")
}
