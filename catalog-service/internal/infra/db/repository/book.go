package repository

type IBook interface {
	ID() int
	Title() string
	Year() int
	AuthorID() int
	Description() *string
	Categories() []int
}
