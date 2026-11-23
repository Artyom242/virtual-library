package command

type CreateAuthorCommand struct {
	Name string
}

type UpdateAuthorCommand struct {
	ID   int
	Name *string
}
