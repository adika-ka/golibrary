package entities

type Book struct {
	ID       int     `db:"id" json:"id"`
	Title    string  `db:"title" json:"title"`
	AuthorID int     `db:"author_id" json:"author_id"`
	Author   *Author `json:"author,omitempty"`
}
