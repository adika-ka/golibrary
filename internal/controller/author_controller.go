package controller

import (
	"encoding/json"
	"fmt"
	"golibrary/infrastructure"
	"golibrary/internal/facade"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type AuthorController struct {
	Facade    *facade.LibraryFacade
	Responder infrastructure.Responder
}

func NewAuthorController(facade *facade.LibraryFacade, responder infrastructure.Responder) *AuthorController {
	return &AuthorController{
		Facade:    facade,
		Responder: responder,
	}
}

func RegisterAuthorRoutes(r chi.Router, ac *AuthorController) {
	r.Route("/authors", func(r chi.Router) {
		r.Get("/", GetAllAuthors(ac))
		r.Post("/", CreateAuthor(ac))
		r.Get("/top", GetTopAuthors(ac))
	})
}

// GetAllAuthors godoc
// @Summary Получить список всех авторов
// @Tags authors
// @Produce json
// @Success 200 {array} entities.Author
// @Router /authors [get]
func GetAllAuthors(ac *AuthorController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authors, err := ac.Facade.GetAllAuthors(r.Context())
		if err != nil {
			ac.Responder.ErrorInternal(w, err)
			return
		}

		ac.Responder.OutputJSON(w, authors)
	}
}

// CreateAuthor godoc
// @Summary Добавить нового автора
// @Tags authors
// @Accept json
// @Produce json
// @Param author body entities.Author true "Имя автора"
// @Success 201 {object} entities.Author
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /authors [post]
func CreateAuthor(ac *AuthorController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var a struct{ Name string }

		if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
			ac.Responder.ErrorBadRequest(w, err)
			return
		}

		if a.Name == "" {
			ac.Responder.ErrorBadRequest(w, fmt.Errorf("author name cannot be empty"))
			return
		}

		author, err := ac.Facade.CreateAuthor(r.Context(), a.Name)
		if err != nil {
			ac.Responder.ErrorInternal(w, err)
			return
		}

		ac.Responder.OutputCreated(w, author)
	}
}

// GetTopAuthors godoc
// @Summary Получить топ авторов по количеству выданных книг
// @Tags authors
// @Produce json
// @Success 200 {array} entities.Author
// @Router /authors/top [get]
func GetTopAuthors(ac *AuthorController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		top, err := ac.Facade.GetTopAuthors(r.Context(), 5)
		if err != nil {
			ac.Responder.ErrorInternal(w, err)
			return
		}

		ac.Responder.OutputJSON(w, top)
	}
}
