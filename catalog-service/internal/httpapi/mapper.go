package httpapi

import (
	"virtual-library/catalog-service/internal/domain"
	"virtual-library/catalog-service/internal/httpapi/dto"
)

func toBookResponse(book domain.Book) dto.BookResponse {
	return dto.BookResponse{
		ID:       book.ID,
		Title:    book.Title,
		Year:     book.Year,
		AuthorID: book.AuthorID,
	}
}

func toAuthorResponse(author domain.Author) dto.AuthorResponse {
	return dto.AuthorResponse{
		ID:   author.ID,
		Name: author.Name,
	}
}
