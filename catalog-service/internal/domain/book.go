package domain

type Book struct {
	ID       int
	Title    string
	AuthorID int
	Year     int
}

type BookPatch struct {
	Title    *string
	AuthorID *int
	Year     *int
}
