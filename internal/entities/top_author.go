package entities

type TopAuthor struct {
	ID               int    `db:"id"`
	Name             string `db:"name"`
	RentedBooksCount int    `db:"rented_books_count"`
}
