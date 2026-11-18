package dto

type BookResponse struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	AuthorID int    `json:"author_id"`
	Year     int    `json:"year"`
}

type BooksResponse struct {
	Data []BookResponse `json:"data"`
}

type BookDataResponse struct {
	Data BookResponse `json:"data"`
}
