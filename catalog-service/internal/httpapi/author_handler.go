package httpapi

import (
	"errors"
	"log/slog"
	"net/http"
	"virtual-library/catalog-service/internal/domain"
	"virtual-library/catalog-service/internal/httpapi/generated"
	"virtual-library/catalog-service/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthorHandler struct {
	service *service.AuthorService
}

func NewAuthorHandler(service *service.AuthorService) *AuthorHandler {
	return &AuthorHandler{
		service: service,
	}
}

func (h *AuthorHandler) GetAuthors(c *gin.Context, params generated.GetAuthorsParams) {
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

	authors, total, err := h.service.GetAll(ctx, page, limit)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get all authors",
			slog.Any("error", err),
		)
		c.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: domain.ErrInternalServer.Error()})
		return
	}

	response := make([]generated.AuthorResponse, 0, len(authors))
	for _, author := range authors {
		response = append(response, toGeneratedAuthorResponse(author))
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
		slog.WarnContext(ctx, "invalid request body",
			slog.Any("error", err),
		)
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: domain.ErrInvalidAuthorData.Error()})
		return
	}

	author, err := domain.NewAuthor(request.Name)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidAuthorData) {
			slog.WarnContext(ctx, "invalid request body",
				slog.Any("error", err),
			)
			c.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: domain.ErrInvalidAuthorData.Error()})
			return
		}
		slog.ErrorContext(ctx, "server internal error",
			slog.Any("error", err),
		)
		c.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: domain.ErrInternalServer.Error()})
		return
	}

	createdAuthor, err := h.service.Create(ctx, author)
	if err != nil {
		if errors.Is(err, domain.ErrAuthorAlreadyExists) {
			slog.WarnContext(ctx, "author already exists",
				slog.String("author_name", author.Name()),
			)
			c.JSON(http.StatusConflict, generated.ErrorResponse{Message: domain.ErrAuthorAlreadyExists.Error()})
			return
		}
		if errors.Is(err, domain.ErrMapAuthorToDomain) {
			slog.ErrorContext(ctx, "failed to map author to domain",
				slog.String("author_name", author.Name()),
			)
			c.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: domain.ErrMapAuthorToDomain.Error()})
			return
		}

		slog.ErrorContext(ctx, "failed to create author",
			slog.Any("error", err),
		)
		c.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: domain.ErrInternalServer.Error()})
		return
	}

	slog.InfoContext(ctx, "created author",
		slog.Group("author",
			slog.Int("id", createdAuthor.ID()),
			slog.String("name", createdAuthor.Name()),
		),
	)
	c.JSON(http.StatusCreated, generated.CreateAuthor201JSONResponse{Data: toGeneratedAuthorResponse(createdAuthor)})
}

func (h *AuthorHandler) GetAuthorById(c *gin.Context, id generated.ID) {
	ctx := c.Request.Context()
	if id < 1 {
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: domain.ErrInvalidAuthorID.Error(),
		})
		return
	}

	author, err := h.service.GetByID(ctx, id)

	if err != nil {
		if errors.Is(err, domain.ErrAuthorNotFound) {
			slog.WarnContext(ctx, "author not found",
				slog.Int("author_id", id),
			)
			c.JSON(http.StatusNotFound, generated.ErrorResponse{Message: domain.ErrAuthorNotFound.Error()})
			return
		}

		slog.ErrorContext(ctx, "failed to get author by id",
			slog.Int("author_id", id),
			slog.Any("error", err),
		)
		c.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: domain.ErrInternalServer.Error()})
		return
	}

	slog.InfoContext(ctx, "author found",
		slog.Group("author",
			slog.Int("id", author.ID()),
			slog.String("name", author.Name()),
		),
	)
	c.JSON(http.StatusOK, generated.GetAuthorById200JSONResponse{Data: toGeneratedAuthorResponse(*author)})
}

func (h *AuthorHandler) UpdateAuthor(c *gin.Context, id generated.ID) {
	ctx := c.Request.Context()
	if id < 1 {
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: domain.ErrInvalidAuthorID.Error(),
		})
		return
	}

	var request generated.UpdateAuthorJSONRequestBody
	if err := c.ShouldBindJSON(&request); err != nil {
		slog.WarnContext(ctx, "invalid request body",
			slog.Any("error", err),
		)
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: domain.ErrInvalidAuthorData.Error()})
		return
	}

	patch := domain.AuthorPatch{
		Name: request.Name,
	}
	author, err := h.service.Update(ctx, id, patch)
	if errors.Is(err, domain.ErrInvalidAuthorData) {
		slog.WarnContext(ctx, "invalid author data",
			slog.Int("author_id", id),
		)
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: domain.ErrInvalidAuthorData.Error()})
		return
	}
	if errors.Is(err, domain.ErrEmptyAuthorPatch) {
		slog.WarnContext(ctx, "empty author patch",
			slog.Int("author_id", id),
		)
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: domain.ErrEmptyAuthorPatch.Error()})
		return
	}
	if errors.Is(err, domain.ErrAuthorNotFound) {
		slog.WarnContext(ctx, "author not found",
			slog.Int("author_id", id),
		)
		c.JSON(http.StatusNotFound, generated.ErrorResponse{Message: domain.ErrAuthorNotFound.Error()})
		return
	}
	if errors.Is(err, domain.ErrAuthorAlreadyExists) {
		slog.ErrorContext(ctx, "author already exists",
			slog.String("author_name", author.Name()),
		)
		c.JSON(http.StatusConflict, generated.ErrorResponse{Message: domain.ErrAuthorAlreadyExists.Error()})
		return
	}
	if err != nil {
		slog.ErrorContext(ctx, "failed to update author",
			slog.Int("author_id", id),
			slog.Any("error", err),
		)
		c.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: domain.ErrInternalServer.Error()})
		return
	}

	slog.InfoContext(ctx, "updated author",
		slog.Group("author",
			slog.Int("id", id),
			slog.String("name", author.Name()),
		),
	)
	c.JSON(http.StatusOK, generated.UpdateAuthor200JSONResponse{Data: toGeneratedAuthorResponse(author)})
}

func (h *AuthorHandler) DeleteAuthor(c *gin.Context, id generated.ID) {
	ctx := c.Request.Context()
	if id < 1 {
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: domain.ErrInvalidAuthorID.Error(),
		})
		return
	}

	err := h.service.Delete(ctx, id)
	if errors.Is(err, domain.ErrAuthorNotFound) {
		slog.WarnContext(ctx, "author not found",
			slog.Int("author_id", id),
		)
		c.JSON(http.StatusNotFound, generated.ErrorResponse{Message: domain.ErrAuthorNotFound.Error()})
		return
	}
	if errors.Is(err, domain.ErrAuthorHasBooks) {
		slog.WarnContext(ctx, "author has books, cannot delete",
			slog.Int("author_id", id),
		)
		c.JSON(http.StatusConflict, generated.ErrorResponse{Message: domain.ErrAuthorHasBooks.Error()})
		return
	}
	if err != nil {
		slog.ErrorContext(ctx, "failed to delete author",
			slog.Int("author_id", id),
			slog.Any("error", err),
		)
		c.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: domain.ErrInternalServer.Error()})
		return
	}

	slog.InfoContext(ctx, "author deleted",
		slog.Int("author_id", id),
	)
	c.Status(http.StatusNoContent)
}
