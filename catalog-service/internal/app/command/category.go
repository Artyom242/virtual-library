package command

type CreateCategoryCommand struct {
	Name string
}

type UpdateCategoryCommand struct {
	ID   int
	Name *string
}
