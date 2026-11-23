package command

type CreateBookCommand struct {
	Title       string
	Year        int
	AuthorID    int
	Description *string
}

type UpdateBookCommand struct {
	ID          int
	Title       *string
	AuthorID    *int
	Year        *int
	Description *string
}
