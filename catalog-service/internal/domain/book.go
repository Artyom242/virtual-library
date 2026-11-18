package domain

type Book struct {
	ID       int
	Title    string
	Year     int
	AuthorID int
}

type BookPatch struct {
	Title    *string
	AuthorID *int
	Year     *int
}
