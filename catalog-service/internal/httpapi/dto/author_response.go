package dto

type AuthorResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type AuthorsResponse struct {
	Data []AuthorResponse `json:"data"`
}

type AuthorDataResponse struct {
	Data AuthorResponse `json:"data"`
}
