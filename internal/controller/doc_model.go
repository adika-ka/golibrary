package controller

// UserResponse представляет пользователя с арендованными книгами.
// @name UserResponse
type UserResponse struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Email       string         `json:"email"`
	RentedBooks []BookResponse `json:"rented_books"`
}

// BookResponse представляет книгу с информацией об авторе.
// @name BookResponse
type BookResponse struct {
	ID     int           `json:"id"`
	Title  string        `json:"title"`
	Author AuthorSummary `json:"author"`
}

// AuthorSummary представляет краткую информацию об авторе.
// @name AuthorSummary
type AuthorSummary struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type BookCreateRequest struct {
	Title    string `json:"title" example:"Some Book Title"`
	AuthorID int    `json:"author_id" example:"52"`
}
