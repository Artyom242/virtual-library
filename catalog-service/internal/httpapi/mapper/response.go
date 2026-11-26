package mapper

import (
	"virtual-library/catalog-service/internal/domain/entity"
	"virtual-library/catalog-service/internal/httpapi/generated"
)

func ToGeneratedBookResponse(book *entity.Book) generated.BookResponse {
	return generated.BookResponse{
		Id:          book.ID(),
		Title:       book.Title(),
		AuthorId:    book.AuthorID(),
		Year:        book.Year(),
		Description: book.Description(),
		Categories:  book.Categories(),
	}
}
func ToGeneratedBookResponses(books []*entity.Book) []generated.BookResponse {
	response := make([]generated.BookResponse, 0, len(books))

	for _, book := range books {
		response = append(response, ToGeneratedBookResponse(book))
	}

	return response
}

func ToGeneratedAuthorResponse(author *entity.Author) generated.AuthorResponse {
	return generated.AuthorResponse{
		Id:   author.ID(),
		Name: author.Name(),
	}
}

func ToGeneratedCategoryResponse(category *entity.Category) generated.CategoryResponse {
	return generated.CategoryResponse{
		Id:   category.ID(),
		Name: category.Name(),
	}
}
