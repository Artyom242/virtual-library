package httpapi

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"virtual-library/catalog-service/internal/domain"
	"virtual-library/catalog-service/internal/httpapi/dto"
	"virtual-library/catalog-service/internal/service"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	service *service.BookService
}

func NewBookHandler(service *service.BookService) *BookHandler {
	return &BookHandler{service: service}
}

func (h *BookHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	data := c.Param("id")
	id, err := strconv.Atoi(data)
	if err != nil || id < 1 {
		slog.WarnContext(ctx, "invalid id",
			slog.String("book_id", data),
			slog.Any("error", err),
		)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: domain.ErrInvalidBookID.Error()})
		return
	}

	book, err := h.service.GetByID(ctx, id)
	if errors.Is(err, domain.ErrBookNotFound) {
		slog.WarnContext(ctx, "book not found",
			slog.Int("book_id", id),
		)
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Message: domain.ErrBookNotFound.Error()})
		return
	}
	if err != nil {
		slog.ErrorContext(ctx, "failed to get book by id",
			slog.Int("book_id", id),
			slog.Any("error", err),
		)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: domain.ErrInternalServer.Error()})
		return
	}
	slog.InfoContext(ctx, "book found",
		slog.Int("book_id", id),
		slog.Any("book", book),
	)
	c.JSON(http.StatusOK, dto.BookDataResponse{
		Data: toBookResponse(*book),
	})
}
func (h *BookHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	books, err := h.service.GetAll(ctx)

	if err != nil {
		slog.ErrorContext(ctx, "failed to get all books",
			slog.Any("error", err),
		)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: domain.ErrInternalServer.Error()})
		return
	}

	response := make([]dto.BookResponse, 0, len(books))

	for _, book := range books {
		response = append(response, toBookResponse(book))
	}

	slog.InfoContext(ctx, "got all books",
		slog.Int("count", len(books)),
	)
	c.JSON(http.StatusOK, dto.BooksResponse{Data: response})
}
func (h *BookHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var request dto.CreateBookRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		slog.WarnContext(ctx, "invalid request body",
			slog.Any("error", err),
		)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: domain.ErrInvalidBookData.Error()})
		return
	}

	book := domain.Book{
		Title:    request.Title,
		AuthorID: request.AuthorID,
		Year:     request.Year,
	}

	createdBook, err := h.service.Create(ctx, book)
	if errors.Is(err, domain.ErrAuthorNotFound) {
		slog.WarnContext(ctx, "author not found for creating book",
			slog.Int("author_id", book.AuthorID),
		)
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Message: domain.ErrAuthorNotFound.Error()})
		return
	}

	if err != nil {
		slog.ErrorContext(ctx, "failed to create book",
			slog.Any("error", err),
		)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: domain.ErrInternalServer.Error()})
		return
	}

	slog.InfoContext(ctx, "created book",
		slog.Group("book",
			slog.Int("id", createdBook.ID),
			slog.Int("year", createdBook.Year),
			slog.String("title", createdBook.Title),
			slog.Int("author_id", createdBook.AuthorID),
		),
	)
	c.JSON(http.StatusCreated, dto.BookDataResponse{
		Data: toBookResponse(createdBook),
	})
}
func (h *BookHandler) CreateBookWithNewAuthor(c *gin.Context) {
	ctx := c.Request.Context()
	var request dto.CreateBookWithAuthorRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		slog.WarnContext(ctx, "invalid request body",
			slog.Any("error", err),
		)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: domain.ErrInvalidBookData.Error()})
		return
	}

	book := domain.Book{
		Title: request.Title,
		Year:  request.Year,
	}

	author := domain.Author{
		Name: request.AuthorName,
	}

	createdBook, err := h.service.CreateBookWithNewAuthor(ctx, book, author)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create book with new author",
			slog.Any("error", err),
		)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: domain.ErrInternalServer.Error()})
		return
	}

	slog.InfoContext(ctx, "created book with new author",
		slog.Group("book",
			slog.Int("id", createdBook.ID),
			slog.Int("year", createdBook.Year),
			slog.String("title", createdBook.Title),
			slog.Int("author_id", createdBook.AuthorID),
			slog.Group("author",
				slog.String("name", author.Name),
			),
		),
	)
	c.JSON(http.StatusCreated, dto.BookDataResponse{
		Data: toBookResponse(createdBook),
	})

}
func (h *BookHandler) DeleteBook(c *gin.Context) {
	ctx := c.Request.Context()
	data := c.Param("id")
	id, err := strconv.Atoi(data)
	if err != nil || id < 1 {
		slog.WarnContext(ctx, "invalid id",
			slog.String("book_id", data),
			slog.Any("error", err),
		)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: domain.ErrInvalidBookID.Error()})
		return
	}

	err = h.service.Delete(ctx, id)
	if errors.Is(err, domain.ErrBookNotFound) {
		slog.WarnContext(ctx, "book not found",
			slog.Int("book_id", id),
		)
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Message: domain.ErrBookNotFound.Error()})
		return
	}

	if err != nil {
		slog.ErrorContext(ctx, "failed to delete book",
			slog.Int("book_id", id),
			slog.Any("error", err),
		)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: domain.ErrInternalServer.Error()})
		return
	}

	slog.InfoContext(ctx, "book deleted",
		slog.Int("book_id", id),
	)
	c.Status(http.StatusNoContent)
}
func (h *BookHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()
	data := c.Param("id")
	id, err := strconv.Atoi(data)
	if err != nil || id < 1 {
		slog.WarnContext(ctx, "invalid id",
			slog.String("book_id", data),
			slog.Any("error", err),
		)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: domain.ErrInvalidBookID.Error()})
		return
	}

	var request dto.UpdateBookRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		slog.ErrorContext(ctx, "invalid request body",
			slog.Any("error", err),
		)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: domain.ErrInvalidBookData.Error()})
		return
	}

	patchBook := domain.BookPatch{
		Title:    request.Title,
		Year:     request.Year,
		AuthorID: request.AuthorID,
	}

	book, err := h.service.Update(ctx, id, patchBook)
	if errors.Is(err, domain.ErrEmptyBookPatch) {
		slog.WarnContext(ctx, "empty book patch",
			slog.Int("book_id", id),
		)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: domain.ErrEmptyBookPatch.Error()})
		return
	}
	if errors.Is(err, domain.ErrBookNotFound) {
		slog.WarnContext(ctx, "book not found",
			slog.Int("book_id", id),
		)
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Message: domain.ErrBookNotFound.Error()})
		return
	}
	if err != nil {
		slog.ErrorContext(ctx, "failed to update book",
			slog.Int("book_id", id),
			slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: domain.ErrInternalServer.Error()})
		return
	}

	slog.InfoContext(ctx, "updated book",
		slog.Group("new_book",
			slog.Int("book_id", id),
			slog.Int("year", book.Year),
			slog.String("title", book.Title),
			slog.Int("author_id", book.AuthorID),
		),
	)
	c.JSON(http.StatusOK, dto.BookDataResponse{Data: toBookResponse(book)})
}
