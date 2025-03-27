package controller

import (
	"encoding/json"
	"fmt"
	"golibrary/infrastructure"
	"golibrary/internal/entities"
	"golibrary/internal/facade"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type BookController struct {
	Facade    *facade.LibraryFacade
	Responder infrastructure.Responder
}

func NewBookController(facade *facade.LibraryFacade, responder infrastructure.Responder) *BookController {
	return &BookController{
		Facade:    facade,
		Responder: responder,
	}
}

func RegisterBookRoutes(r chi.Router, bc *BookController) {
	r.Route("/books", func(r chi.Router) {
		r.Get("/", GetBooks(bc))
		r.Post("/", CreateBook(bc))
	})
}

// GetBooks godoc
// @Summary Получить список всех книг
// @Tags books
// @Produce json
// @Success 200 {array} entities.Book
// @Router /books [get]
func GetBooks(bc *BookController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		books, err := bc.Facade.GetAllBooks(r.Context())
		if err != nil {
			bc.Responder.ErrorInternal(w, err)
			return
		}

		bc.Responder.OutputJSON(w, books)
	}
}

// CreateBook godoc
// @Summary Добавить новую книгу
// @Description Добавляет книгу с указанным названием и ID автора
// @Tags books
// @Accept json
// @Produce json
// @Param book body BookCreateRequest true "Данные книги"
// @Success 201 {object} entities.Book
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books [post]
func CreateBook(bc *BookController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b entities.Book

		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			bc.Responder.ErrorBadRequest(w, err)
			return
		}

		book, err := bc.Facade.CreateBook(r.Context(), b)
		if err != nil {
			bc.Responder.ErrorInternal(w, err)
			return
		}

		bc.Responder.OutputCreated(w, book)
	}
}

func RentBook(bc *BookController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.URL.Query().Get("user_id")

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			bc.Responder.ErrorBadRequest(w, fmt.Errorf("invalid user ID"))
			return
		}

		bookIDStr := r.URL.Query().Get("book_id")

		bookID, err := strconv.Atoi(bookIDStr)
		if err != nil {
			bc.Responder.ErrorBadRequest(w, fmt.Errorf("invalid book ID"))
		}

		err = bc.Facade.RentBook(userID, bookID)
		if err != nil {
			bc.Responder.ErrorInternal(w, err)
			return
		}

		bc.Responder.OutputNoContent(w)
	}
}

func ReturnBook(bc *BookController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.URL.Query().Get("user_id")

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			bc.Responder.ErrorBadRequest(w, fmt.Errorf("invalid user ID"))
			return
		}

		bookIDStr := r.URL.Query().Get("book_id")

		bookID, err := strconv.Atoi(bookIDStr)
		if err != nil {
			bc.Responder.ErrorBadRequest(w, fmt.Errorf("invalid book ID"))
		}

		err = bc.Facade.ReturnBook(userID, bookID)
		if err != nil {
			bc.Responder.ErrorInternal(w, err)
			return
		}

		bc.Responder.OutputNoContent(w)
	}
}
