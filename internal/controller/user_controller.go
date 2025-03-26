package controller

import (
	"context"
	"golibrary/infrastructure"
	"golibrary/internal/facade"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type UserController struct {
	Facade    *facade.LibraryFacade
	Responder infrastructure.Responder
}

func NewUserController(facade *facade.LibraryFacade, responder infrastructure.Responder) *UserController {
	return &UserController{
		Facade:    facade,
		Responder: responder,
	}
}

func RegisterUserRoutes(r chi.Router, uc *UserController) {
	r.Route("/users", func(r chi.Router) {
		r.Get("/", GetAllUsers(uc))
		r.Get("/borrowed-books", uc.GetUserBorrowedBooksEndpoint)
	})
}

// GetAllUsers godoc
// @Summary Получить список всех пользователей
// @Tags users
// @Produce json
// @Success 200 {array} entities.User
// @Router /users [get]
func GetAllUsers(uc *UserController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := uc.Facade.GetAllUsers(r.Context())
		if err != nil {
			uc.Responder.ErrorInternal(w, err)
			return
		}

		uc.Responder.OutputJSON(w, users)
	}
}

// GetUserBorrowedBooksEndpoint godoc
// @Summary Получить список книг, взятых пользователем в аренду (активные займы)
// @Tags users
// @Produce json
// @Param user_id query int true "ID пользователя"
// @Success 200 {array} entities.Book "Список книг, которые пользователь не вернул"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Router /users/borrowed-books [get]
func (uc *UserController) GetUserBorrowedBooksEndpoint(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		uc.Responder.ErrorBadRequest(w, err)
		return
	}

	books, err := uc.Facade.GetBooksBorrowedByUser(context.Background(), userID)
	if err != nil {
		uc.Responder.ErrorInternal(w, err)
		return
	}

	uc.Responder.OutputJSON(w, books)
}
