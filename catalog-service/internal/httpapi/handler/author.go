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

type AuthorHandler struct {
	service interfaces.AuthorService
}

func NewAuthorHandler(service interfaces.AuthorService) *AuthorHandler {
	return &AuthorHandler{
		service: service,
	}
}

func (h *AuthorHandler) GetAuthors(c *gin.Context, params generated.GetAuthorsParams) {
	ctx := c.Request.Context()

	paginationParams := pagination.NewPagination(params.Page, params.Limit)

	authors, total, err := h.service.GetAllAuthors(ctx, paginationParams)
	if err != nil {
		errorHandler.HandleCatalogError(ctx, c, err)
		return
	}

	response := make([]generated.AuthorResponse, 0, len(authors))
	for _, author := range authors {
		response = append(response, mapper.ToGeneratedAuthorResponse(author))
	}

	slog.InfoContext(ctx, "got all authors",
		slog.Int("count", len(authors)),
		slog.Int("total", total),
	)
	c.JSON(http.StatusOK, generated.GetAuthors200JSONResponse{
		Data:  response,
		Total: total,
	})
}

func (h *AuthorHandler) CreateAuthor(c *gin.Context) {
	ctx := c.Request.Context()
	var request generated.CreateAuthorJSONRequestBody

	if err := c.ShouldBindJSON(&request); err != nil {
		errorHandler.HandleCatalogError(ctx, c, err, slog.String("author_name", request.Name))
		return
	}

	authorCommand := command.CreateAuthorCommand{
		Name: request.Name,
	}

	createdAuthor, err := h.service.CreateAuthor(ctx, &authorCommand)
	if err != nil {
		errorHandler.HandleCatalogError(ctx, c, err, slog.String("author_name", authorCommand.Name))
		return
	}

	slog.InfoContext(ctx, "created author",
		slog.Group("author",
			slog.Int("id", createdAuthor.ID()),
			slog.String("name", createdAuthor.Name()),
		),
	)
	c.JSON(http.StatusOK, generated.CreateAuthor201JSONResponse{
		Data: mapper.ToGeneratedAuthorResponse(createdAuthor),
	})
}

func (h *AuthorHandler) GetAuthorById(c *gin.Context, id generated.ID) {
	ctx := c.Request.Context()
	if !validation.IsIDMoreThanOne(id) {
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: entities.ErrInvalidAuthorID.Error(),
		})
		return
	}

	author, err := h.service.GetAuthorByID(ctx, id)
	if err != nil {
		errorHandler.HandleCatalogError(ctx, c, err, slog.Int("author_id", id))
		return
	}

	slog.InfoContext(ctx, "author found",
		slog.Group("author",
			slog.Int("id", author.ID()),
			slog.String("name", author.Name()),
		),
	)
	c.JSON(http.StatusOK, generated.GetAuthorById200JSONResponse{
		Data: mapper.ToGeneratedAuthorResponse(author),
	})
}

func (h *AuthorHandler) UpdateAuthor(c *gin.Context, id generated.ID) {
	ctx := c.Request.Context()
	if !validation.IsIDMoreThanOne(id) {
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: entities.ErrInvalidAuthorID.Error(),
		})
		return
	}

	var request generated.UpdateAuthorJSONRequestBody
	if err := c.ShouldBindJSON(&request); err != nil {
		errorHandler.HandleCatalogError(ctx, c, err, slog.Int("author_id", id))
		return
	}

	authorCommand := command.UpdateAuthorCommand{
		ID:   id,
		Name: request.Name,
	}

	author, err := h.service.UpdateAuthor(ctx, &authorCommand)
	if err != nil {
		attrs := []slog.Attr{slog.Int("author_id", authorCommand.ID)}
		if request.Name != nil {
			attrs = append(attrs, slog.String("author_name", *authorCommand.Name))
		}
		errorHandler.HandleCatalogError(ctx, c, err, attrs...)
		return
	}

	slog.InfoContext(ctx, "updated author",
		slog.Group("author",
			slog.Int("id", id),
			slog.String("name", author.Name()),
		),
	)
	c.JSON(http.StatusOK, generated.UpdateAuthor200JSONResponse{
		Data: mapper.ToGeneratedAuthorResponse(author),
	})
}

func (h *AuthorHandler) DeleteAuthor(c *gin.Context, id generated.ID) {
	ctx := c.Request.Context()
	if !validation.IsIDMoreThanOne(id) {
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: entities.ErrInvalidAuthorID.Error(),
		})
		return
	}

	err := h.service.DeleteAuthor(ctx, id)
	if err != nil {
		errorHandler.HandleCatalogError(ctx, c, err, slog.Int("author_id", id))
		return
	}

	slog.InfoContext(ctx, "author deleted",
		slog.Int("author_id", id),
	)
	c.Status(http.StatusNoContent)
}
