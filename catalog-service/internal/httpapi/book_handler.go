package httpapi

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"virtual-library/catalog-service/internal/domain"
	"virtual-library/catalog-service/internal/httpapi/generated"
	"virtual-library/catalog-service/internal/service"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	service *service.BookService
}

type BookLogContext struct {
	BookID   *int
	AuthorID *int
}

func NewBookHandler(service *service.BookService) *BookHandler {
	return &BookHandler{service: service}
}

func (h *BookHandler) GetBooks(c *gin.Context, params generated.GetBooksParams) {
	ctx := c.Request.Context()

	page, limit, err := parsePagination(params.Page, params.Limit)
	if err != nil {
		slog.WarnContext(ctx, "invalid pagination parameters",
			slog.Any("error", err),
		)
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: domain.ErrInvalidPagination.Error(),
		})
		return
	}

	books, total, err := h.service.GetAll(ctx, page, limit)
	if err != nil {
		handleBookError(ctx, c, err, BookLogContext{})
		return
	}

	response := make([]generated.BookResponse, 0, len(books))

	for _, book := range books {
		response = append(response, toGeneratedBookResponse(book))
	}

	slog.InfoContext(ctx, "got all books",
		slog.Int("count", len(books)),
	)
	c.JSON(http.StatusOK, generated.GetBooks200JSONResponse{
		Data:  response,
		Total: total,
	})
}

func (h *BookHandler) GetBookById(c *gin.Context, id generated.ID) {
	ctx := c.Request.Context()
	if id < 1 {
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: domain.ErrInvalidBookID.Error(),
		})
		return
	}

	book, err := h.service.GetByID(ctx, id)
	if err != nil {
		logCtx := BookLogContext{
			BookID: &id,
		}
		handleBookError(ctx, c, err, logCtx)
		return
	}
	slog.InfoContext(ctx, "book found",
		slog.Int("book_id", id),
		slog.Any("book", book),
	)
	c.JSON(http.StatusOK, generated.GetBookById200JSONResponse{
		Data: toGeneratedBookResponse(*book),
	})
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	ctx := c.Request.Context()
	var request generated.CreateBookJSONRequestBody

	if err := c.ShouldBindJSON(&request); err != nil {
		slog.WarnContext(ctx, "invalid request body",
			slog.Any("error", err),
		)
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: domain.ErrInvalidBookData.Error()})
		return
	}
	if request.AuthorId < 1 {
		slog.WarnContext(ctx, "invalid request body",
			slog.Int("author_id", request.AuthorId),
		)
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: domain.ErrInvalidAuthorID.Error()})
		return
	}

	book, err := domain.NewBook(
		request.Title,
		request.Year,
		request.AuthorId,
		request.Description,
	)
	if err != nil {
		slog.WarnContext(ctx, "invalid book data",
			slog.Any("error", err),
		)
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: domain.ErrInvalidBookData.Error(),
		})
		return
	}

	createdBook, err := h.service.Create(ctx, book)
	if err != nil {
		logCtx := BookLogContext{
			AuthorID: &request.AuthorId,
		}
		handleBookError(ctx, c, err, logCtx)
		return
	}

	response := toGeneratedBookResponse(createdBook)
	slog.InfoContext(ctx, "created book",
		slog.Group("book",
			slog.Int("id", createdBook.ID()),
			slog.Int("year", createdBook.Year()),
			slog.String("title", createdBook.Title()),
			slog.Int("author_id", createdBook.AuthorID()),
			slog.Any("description", createdBook.Description()),
		),
	)
	c.JSON(http.StatusCreated, generated.CreateBook201JSONResponse{
		Data: response,
	})
}

func (h *BookHandler) CreateBookWithAuthor(c *gin.Context) {
	ctx := c.Request.Context()
	var request generated.CreateBookWithAuthorJSONRequestBody

	if err := c.ShouldBindJSON(&request); err != nil {
		slog.WarnContext(ctx, "invalid request body",
			slog.Any("error", err),
		)
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: domain.ErrInvalidBookData.Error()})
		return
	}

	if request.Year < 1 {
		slog.WarnContext(ctx, "invalid request year",
			slog.Any("year", request.Year),
		)
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: domain.ErrInvalidBookData.Error(),
		})
		return
	}

	author, err := domain.NewAuthor(request.AuthorName)
	if err != nil {
		slog.WarnContext(ctx, "invalid author data",
			slog.Any("error", err),
		)
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: domain.ErrInvalidAuthorData.Error()})
		return
	}

	book := domain.CreateBookWithAuthor{
		Title:       request.Title,
		Year:        request.Year,
		Description: request.Description,
	}

	createdBook, err := h.service.CreateBookWithNewAuthor(ctx, book, author)
	if err != nil {
		handleBookError(ctx, c, err, BookLogContext{})
		return
	}

	response := toGeneratedBookResponse(createdBook)
	slog.InfoContext(ctx, "created book with new author",
		slog.Group("book",
			slog.Int("id", createdBook.ID()),
			slog.Int("year", createdBook.Year()),
			slog.String("title", createdBook.Title()),
			slog.Int("author_id", createdBook.AuthorID()),
			slog.Any("description", createdBook.Description()),
			slog.Group("author",
				slog.String("name", author.Name()),
			),
		),
	)
	c.JSON(http.StatusCreated, generated.CreateBookWithAuthor201JSONResponse{
		Data: response,
	})
}

func (h *BookHandler) DeleteBook(c *gin.Context, id generated.ID) {
	ctx := c.Request.Context()
	if id < 1 {
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: domain.ErrInvalidBookID.Error(),
		})
		return
	}

	err := h.service.Delete(ctx, id)
	if err != nil {
		logCtx := BookLogContext{
			BookID: &id,
		}
		handleBookError(ctx, c, err, logCtx)
		return
	}

	slog.InfoContext(ctx, "book deleted",
		slog.Int("book_id", id),
	)
	c.Status(http.StatusNoContent)
}

func (h *BookHandler) UpdateBook(c *gin.Context, id generated.ID) {
	ctx := c.Request.Context()
	if id < 1 {
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: domain.ErrInvalidBookID.Error(),
		})
		return
	}

	var request generated.UpdateBookJSONRequestBody
	if err := c.ShouldBindJSON(&request); err != nil {
		slog.ErrorContext(ctx, "invalid request body",
			slog.Any("error", err),
		)
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: domain.ErrInvalidBookData.Error()})
		return
	}

	patchBook := domain.BookPatch{
		Title:       request.Title,
		Year:        request.Year,
		AuthorID:    request.AuthorId,
		Description: request.Description,
	}
	book, err := h.service.Update(ctx, id, patchBook)
	if err != nil {
		logCtx := BookLogContext{
			BookID:   &id,
			AuthorID: request.AuthorId,
		}
		handleBookError(ctx, c, err, logCtx)
		return
	}

	slog.InfoContext(ctx, "updated book",
		slog.Group("new_book",
			slog.Int("book_id", id),
			slog.Int("year", book.Year()),
			slog.String("title", book.Title()),
			slog.Int("author_id", book.AuthorID()),
			slog.String("description", *book.Description()),
		),
	)
	c.JSON(http.StatusOK, generated.UpdateBook200JSONResponse{Data: toGeneratedBookResponse(book)})
}

func handleBookError(
	ctx context.Context,
	c *gin.Context,
	err error,
	logCtx BookLogContext,
) {
	switch {
	case errors.Is(err, domain.ErrBookNotFound):
		attrs := requestLogAttrs(c)
		attrs = append(attrs, slog.Int("book_id", *logCtx.BookID))

		slog.LogAttrs(ctx, slog.LevelWarn, "book not found", attrs...)
		c.JSON(http.StatusNotFound, generated.ErrorResponse{
			Message: domain.ErrBookNotFound.Error(),
		})

	case errors.Is(err, domain.ErrAuthorNotFound):
		attrs := requestLogAttrs(c)
		attrs = append(attrs, slog.Int("author_id", *logCtx.AuthorID))

		slog.LogAttrs(ctx, slog.LevelWarn, "author not found", attrs...)
		c.JSON(http.StatusNotFound, generated.ErrorResponse{
			Message: domain.ErrAuthorNotFound.Error(),
		})

	case errors.Is(err, domain.ErrEmptyBookPatch):
		attrs := requestLogAttrs(c)
		attrs = append(attrs, slog.Int("book_id", *logCtx.BookID))

		slog.LogAttrs(ctx, slog.LevelWarn, "empty book patch", attrs...)
		c.JSON(http.StatusNotFound, generated.ErrorResponse{
			Message: domain.ErrEmptyBookPatch.Error(),
		})

	case errors.Is(err, domain.ErrInvalidBookDescription):
		attrs := requestLogAttrs(c)
		attrs = append(attrs, slog.Int("book_id", *logCtx.BookID))

		slog.LogAttrs(ctx, slog.LevelWarn, "invalid book description", attrs...)
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: domain.ErrInvalidBookDescription.Error(),
		})

	default:
		attrs := requestLogAttrs(c)
		attrs = append(attrs, slog.String("error", err.Error()))

		slog.LogAttrs(ctx, slog.LevelError, "book operation failed", attrs...)
		c.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: domain.ErrInternalServer.Error(),
		})
	}
}
