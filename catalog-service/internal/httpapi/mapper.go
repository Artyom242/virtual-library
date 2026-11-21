package httpapi

import (
	"virtual-library/catalog-service/internal/domain"
	"virtual-library/catalog-service/internal/httpapi/generated"
)

func toGeneratedBookResponse(book domain.Book) generated.BookResponse {
	description := book.Description()

	return generated.BookResponse{
		Id:          book.ID(),
		Title:       book.Title(),
		AuthorId:    book.AuthorID(),
		Year:        book.Year(),
		Description: description,
	}
}

func toGeneratedAuthorResponse(author domain.Author) generated.AuthorResponse {
	return generated.AuthorResponse{
		Id:   author.ID(),
		Name: author.Name(),
	}
}
