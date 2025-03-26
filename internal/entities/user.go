package entities

type User struct {
	ID          int    `db:"id" json:"ID"`
	Name        string `db:"name" json:"Name"`
	Email       string `db:"email" json:"Email"`
	RentedBooks []Book `db:"rented_books" json:"RentedBooks"`
}
