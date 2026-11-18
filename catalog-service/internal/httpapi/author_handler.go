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

type AuthorHandler struct {
	service *service.AuthorService
}

func NewAuthorHandler(service *service.AuthorService) *AuthorHandler {
	return &AuthorHandler{
		service: service,
	}
}

func (h *AuthorHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var request dto.CreateAuthorRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		slog.WarnContext(ctx, "invalid request body",
			slog.Any("error", err),
		)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: domain.ErrInvalidAuthorData.Error()})
		return
	}

	author := domain.Author{
		Name: request.Name,
	}

	createdAuthor, err := h.service.Create(ctx, author)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create author",
			slog.Any("error", err),
		)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: domain.ErrInternalServer.Error()})
		return
	}

	slog.InfoContext(ctx, "created author",
		slog.Group("author",
			slog.Int("id", createdAuthor.ID),
			slog.String("name", createdAuthor.Name),
		),
	)
	c.JSON(http.StatusCreated, dto.AuthorDataResponse{Data: toAuthorResponse(createdAuthor)})
}

func (h *AuthorHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	authors, err := h.service.GetAll(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get all authors",
			slog.Any("error", err),
		)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: domain.ErrInternalServer.Error()})
		return
	}

	response := make([]dto.AuthorResponse, 0, len(authors))
	for _, author := range authors {
		response = append(response, toAuthorResponse(author))
	}

	slog.InfoContext(ctx, "got all authors",
		slog.Int("count", len(authors)),
	)
	c.JSON(http.StatusOK, dto.AuthorsResponse{Data: response})
}

func (h *AuthorHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	request := c.Param("id")
	id, err := strconv.Atoi(request)
	if err != nil || id < 1 {
		slog.WarnContext(ctx, "invalid id",
			slog.String("author_id", request),
			slog.Any("error", err),
		)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: domain.ErrInvalidAuthorID.Error()})
		return
	}

	author, err := h.service.GetByID(ctx, id)

	if errors.Is(err, domain.ErrAuthorNotFound) {
		slog.WarnContext(ctx, "author not found",
			slog.Int("author_id", id),
		)
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Message: domain.ErrAuthorNotFound.Error()})
		return
	}
	if err != nil {
		slog.ErrorContext(ctx, "failed to get author by id",
			slog.Int("author_id", id),
			slog.Any("error", err),
		)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: domain.ErrInternalServer.Error()})
		return
	}

	slog.InfoContext(ctx, "author found",
		slog.Group("author",
			slog.Int("id", author.ID),
			slog.String("name", author.Name),
		),
	)
	c.JSON(http.StatusOK, dto.AuthorDataResponse{Data: toAuthorResponse(*author)})
}

func (h *AuthorHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()
	data := c.Param("id")
	id, err := strconv.Atoi(data)
	if err != nil || id < 1 {
		slog.WarnContext(ctx, "invalid id",
			slog.String("author_id", data),
			slog.Any("error", err),
		)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: domain.ErrInvalidAuthorID.Error()})
		return
	}

	var request dto.UpdateAuthorRequest
	if err = c.ShouldBindJSON(&request); err != nil {
		slog.WarnContext(ctx, "invalid request body",
			slog.Any("error", err),
		)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: domain.ErrInvalidAuthorData.Error()})
		return
	}

	patch := domain.AuthorPatch{
		Name: request.Name,
	}
	author, err := h.service.Update(ctx, id, patch)
	if errors.Is(err, domain.ErrEmptyAuthorPatch) {
		slog.WarnContext(ctx, "empty author patch",
			slog.Int("author_id", id),
		)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: domain.ErrEmptyAuthorPatch.Error()})
		return
	}
	if errors.Is(err, domain.ErrAuthorNotFound) {
		slog.WarnContext(ctx, "author not found",
			slog.Int("author_id", id),
		)
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Message: domain.ErrAuthorNotFound.Error()})
		return
	}
	if err != nil {
		slog.ErrorContext(ctx, "failed to update author",
			slog.Int("author_id", id),
			slog.Any("error", err),
		)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: domain.ErrInternalServer.Error()})
		return
	}

	slog.InfoContext(ctx, "updated author",
		slog.Group("author",
			slog.Int("id", id),
			slog.String("name", author.Name),
		),
	)
	c.JSON(http.StatusOK, dto.AuthorDataResponse{
		Data: toAuthorResponse(author),
	})
}

func (h *AuthorHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	data := c.Param("id")
	id, err := strconv.Atoi(data)
	if err != nil {
		slog.WarnContext(ctx, "invalid id",
			slog.String("author_id", c.Param("id")),
			slog.Any("error", err),
		)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: domain.ErrInvalidAuthorID.Error()})
		return
	}

	err = h.service.Delete(ctx, id)
	if errors.Is(err, domain.ErrAuthorNotFound) {
		slog.WarnContext(ctx, "author not found",
			slog.Int("author_id", id),
		)
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Message: domain.ErrAuthorNotFound.Error()})
		return
	}
	if err != nil {
		slog.ErrorContext(ctx, "failed to delete author",
			slog.Int("author_id", id),
			slog.Any("error", err),
		)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: domain.ErrInternalServer.Error()})
		return
	}

	slog.InfoContext(ctx, "author deleted",
		slog.Int("author_id", id),
	)
	c.Status(http.StatusNoContent)
}
