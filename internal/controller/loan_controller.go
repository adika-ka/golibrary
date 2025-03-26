package controller

import (
	"encoding/json"
	"golibrary/infrastructure"
	"golibrary/internal/facade"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type LoanController struct {
	Facade    *facade.LibraryFacade
	Responder infrastructure.Responder
}

func NewLoanController(facade *facade.LibraryFacade, responder infrastructure.Responder) *LoanController {
	return &LoanController{Facade: facade, Responder: responder}
}

type CreateLoanRequest struct {
	UserID int `json:"user_id"`
	BookID int `json:"book_id"`
}

func RegisterLoanRoutes(r chi.Router, lc *LoanController) {
	r.Route("/loans", func(r chi.Router) {
		r.Post("/borrow", lc.BorrowBookEndpoint)
		r.Post("/return", lc.ReturnBookEndpoint)
	})
}

// BorrowBookEndpoint godoc
// @Summary Выдать книгу пользователю (создать запись аренды)
// @Tags loans
// @Accept json
// @Produce json
// @Param loan body controller.CreateLoanRequest true "Данные аренды (user_id, book_id)"
// @Success 204 "Книга успешно выдана"
// @Failure 400 {object} map[string]string "Неверный запрос или книга уже выдана"
// @Router /loans/borrow [post]
func (lc *LoanController) BorrowBookEndpoint(w http.ResponseWriter, r *http.Request) {
	var req CreateLoanRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lc.Responder.ErrorBadRequest(w, err)
		return
	}

	if err := lc.Facade.RentBook(req.UserID, req.BookID); err != nil {
		lc.Responder.ErrorBadRequest(w, err)
		return
	}

	lc.Responder.OutputNoContent(w)
}

// ReturnBookEndpoint godoc
// @Summary Принять книгу обратно в библиотеку (закрыть активный заем)
// @Tags loans
// @Produce json
// @Param user_id query int true "ID пользователя"
// @Param book_id query int true "ID книги"
// @Success 204 "Книга успешно возвращена"
// @Failure 400 {object} map[string]string "Неверный запрос или книга не числится в аренде"
// @Router /loans/return [post]
func (lc *LoanController) ReturnBookEndpoint(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	bookIDStr := r.URL.Query().Get("book_id")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		lc.Responder.ErrorBadRequest(w, err)
		return
	}

	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil {
		lc.Responder.ErrorBadRequest(w, err)
		return
	}

	if err := lc.Facade.ReturnBook(userID, bookID); err != nil {
		lc.Responder.ErrorBadRequest(w, err)
		return
	}

	lc.Responder.OutputNoContent(w)
}
