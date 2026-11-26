package handler

import (
	"log/slog"
	"net/http"
	"virtual-library/catalog-service/internal/app/command"
	"virtual-library/catalog-service/internal/app/interfaces"
	"virtual-library/catalog-service/internal/domain/entities"
	errorHandler "virtual-library/catalog-service/internal/httpapi/errorhandler"
	"virtual-library/catalog-service/internal/httpapi/generated"
	"virtual-library/catalog-service/internal/httpapi/mapper"
	"virtual-library/catalog-service/internal/httpapi/pagination"
	"virtual-library/catalog-service/internal/httpapi/validation"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	service interfaces.BookService
}

func NewBookHandler(service interfaces.BookService) *BookHandler {
	return &BookHandler{service: service}
}

func (h *BookHandler) GetBooks(c *gin.Context, params generated.GetBooksParams) {
	ctx := c.Request.Context()

	paginationParams := pagination.NewPagination(params.Page, params.Limit)
	books, total, err := h.service.GetAllBooks(ctx, paginationParams)
	if err != nil {
		errorHandler.HandleCatalogError(ctx, c, err)
		return
	}

	response := mapper.ToGeneratedBookResponses(books)

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
	if !validation.IsIDMoreThanOne(id) {
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: entities.ErrInvalidBookID.Error(),
		})
		return
	}

	book, err := h.service.GetBookByID(ctx, id)
	if err != nil {
		errorHandler.HandleCatalogError(ctx, c, err, slog.Int("book_id", id))
		return
	}

	slog.InfoContext(ctx, "book found",
		slog.Int("book_id", id),
		slog.Any("book", book),
	)
	c.JSON(http.StatusOK, generated.GetBookById200JSONResponse{
		Data: mapper.ToGeneratedBookResponse(book),
	})
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	ctx := c.Request.Context()
	var request generated.CreateBookJSONRequestBody

	if err := c.ShouldBindJSON(&request); err != nil {
		errorHandler.HandleCatalogError(ctx, c, err, slog.String("title", request.Title))
		return
	}
	if !validation.IsIDMoreThanOne(request.AuthorId) {
		errorHandler.HandleCatalogError(ctx, c, entities.ErrInvalidAuthorID,
			slog.Int("author_id", request.AuthorId),
		)
		return
	}

	if request.CategoryIds != nil {
		if !validation.IsUniqueValues(*request.CategoryIds) {
			errorHandler.HandleCatalogError(ctx, c, entities.ErrCategoryNotUnique,
				slog.Any("categories_ids", request.CategoryIds),
			)
			return
		}
	}

	book := command.CreateBookCommand{
		Title:       request.Title,
		Year:        request.Year,
		AuthorID:    request.AuthorId,
		Description: request.Description,
		CategoryIDs: request.CategoryIds,
	}

	createdBook, err := h.service.CreateBook(ctx, &book)
	if err != nil {
		errorHandler.HandleCatalogError(ctx, c, err,
			slog.Int("author_id", request.AuthorId),
			slog.String("title", request.Title),
		)
		return
	}

	response := mapper.ToGeneratedBookResponse(createdBook)
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
		errorHandler.HandleCatalogError(ctx, c, err,
			slog.String("title", request.Title),
			slog.String("author_name", request.AuthorName),
		)
		return
	}

	if *request.CategoryIds != nil {
		if !validation.IsUniqueValues(*request.CategoryIds) {
			errorHandler.HandleCatalogError(ctx, c, entities.ErrCategoryNotUnique,
				slog.Any("categories_ids", *request.CategoryIds),
			)
			return
		}
	}

	commandBookWithAuthor := command.CreateBookCommand{
		Title:       request.Title,
		Year:        request.Year,
		Description: request.Description,
		CategoryIDs: request.CategoryIds,
	}
	commandAuthor := command.CreateAuthorCommand{
		Name: request.AuthorName,
	}

	createdBook, err := h.service.CreateBookWithNewAuthor(ctx, &commandBookWithAuthor, &commandAuthor)
	if err != nil {
		errorHandler.HandleCatalogError(ctx, c, err, slog.String("title", request.Title), slog.String("author_name", request.AuthorName))
		return
	}

	response := mapper.ToGeneratedBookResponse(createdBook)
	slog.InfoContext(ctx, "created book with new author",
		slog.Group("book",
			slog.Int("id", createdBook.ID()),
			slog.Int("year", createdBook.Year()),
			slog.String("title", createdBook.Title()),
			slog.Int("author_id", createdBook.AuthorID()),
			slog.Any("description", createdBook.Description()),
		),
	)
	c.JSON(http.StatusCreated, generated.CreateBookWithAuthor201JSONResponse{
		Data: response,
	})
}

func (h *BookHandler) DeleteBook(c *gin.Context, id generated.ID) {
	ctx := c.Request.Context()
	if !validation.IsIDMoreThanOne(id) {
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: entities.ErrInvalidBookID.Error(),
		})
		return
	}

	err := h.service.DeleteBook(ctx, id)
	if err != nil {
		errorHandler.HandleCatalogError(ctx, c, err, slog.Int("book_id", id))
		return
	}

	slog.InfoContext(ctx, "book deleted",
		slog.Int("book_id", id),
	)
	c.Status(http.StatusNoContent)
}

func (h *BookHandler) UpdateBook(c *gin.Context, id generated.ID) {
	ctx := c.Request.Context()
	if !validation.IsIDMoreThanOne(id) {
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: entities.ErrInvalidBookID.Error(),
		})
		return
	}

	var request generated.UpdateBookJSONRequestBody
	if err := c.ShouldBindJSON(&request); err != nil {
		errorHandler.HandleCatalogError(ctx, c, err, slog.Int("book_id", id))
		return
	}

	bookCommand := command.UpdateBookCommand{
		ID:          id,
		Title:       request.Title,
		Year:        request.Year,
		AuthorID:    request.AuthorId,
		Description: request.Description,
		CategoryIDs: request.CategoryIds,
	}

	book, err := h.service.UpdateBook(ctx, &bookCommand)
	if err != nil {
		attrs := []slog.Attr{slog.Int("book_id", id)}
		if request.AuthorId != nil {
			attrs = append(attrs, slog.Int("author_id", *request.AuthorId))
		}
		errorHandler.HandleCatalogError(ctx, c, err, attrs...)
		return
	}

	slog.InfoContext(ctx, "updated book",
		slog.Group("new_book",
			slog.Int("book_id", id),
			slog.Int("year", book.Year()),
			slog.String("title", book.Title()),
			slog.Int("author_id", book.AuthorID()),
		),
	)
	c.JSON(http.StatusOK, generated.UpdateBook200JSONResponse{
		Data: mapper.ToGeneratedBookResponse(book),
	})
}
