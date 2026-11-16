package domain

type Author struct {
	ID   int
	Name string
}

type AuthorPatch struct {
	Name *string
}
