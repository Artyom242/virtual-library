package dto

type CreateBookRequest struct {
	Title    string `json:"title" binding:"required"`
	AuthorID int    `json:"author_id" binding:"required,min=1"`
	Year     int    `json:"year" binding:"required,min=1"`
}
type UpdateBookRequest struct {
	Title    *string `json:"title"`
	AuthorID *int    `json:"author_id"`
	Year     *int    `json:"year"`
}

type CreateBookWithAuthorRequest struct {
	Title string `json:"title" binding:"required"`
	Year  int    `json:"year"  binding:"required"`

	AuthorName string `json:"author_name" binding:"required"`
}
