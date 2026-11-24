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

type CategoryHandler struct {
	service interfaces.CategoryService
}

func NewCategoryHandler(service interfaces.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		service: service,
	}
}

func (h *CategoryHandler) GetCategories(c *gin.Context, params generated.GetCategoriesParams) {
	ctx := c.Request.Context()

	paginationParams := pagination.NewPagination(params.Page, params.Limit)

	categories, total, err := h.service.GetAllCategories(ctx, paginationParams)
	if err != nil {
		errorHandler.HandleCatalogError(ctx, c, err)
		return
	}

	response := make([]generated.CategoryResponse, 0, len(categories))
	for _, category := range categories {
		response = append(response, mapper.ToGeneratedCategoryResponse(category))
	}

	slog.InfoContext(ctx, "got all categories",
		slog.Int("count", len(categories)),
		slog.Int("total", total),
	)
	c.JSON(http.StatusOK, generated.GetCategories200JSONResponse{
		Data:  response,
		Total: total,
	})
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	ctx := c.Request.Context()
	var request generated.CreateCategoryJSONRequestBody

	if err := c.ShouldBindJSON(&request); err != nil {
		errorHandler.HandleCatalogError(ctx, c, err, slog.String("category_name", request.Name))
		return
	}

	categoryCommand := command.CreateCategoryCommand{
		Name: request.Name,
	}

	createdCategory, err := h.service.CreateCategory(ctx, &categoryCommand)
	if err != nil {
		errorHandler.HandleCatalogError(ctx, c, err, slog.String("category_name", categoryCommand.Name))
		return
	}

	slog.InfoContext(ctx, "created category",
		slog.Group("category",
			slog.Int("id", createdCategory.ID()),
			slog.String("name", createdCategory.Name()),
		),
	)
	c.JSON(http.StatusCreated, generated.CreateCategory201JSONResponse{
		Data: mapper.ToGeneratedCategoryResponse(createdCategory),
	})
}

func (h *CategoryHandler) GetCategoryById(c *gin.Context, id generated.ID) {
	ctx := c.Request.Context()
	if !validation.ValidID(id) {
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: entities.ErrInvalidCategoryID.Error(),
		})
		return
	}

	category, err := h.service.GetCategoryByID(ctx, id)
	if err != nil {
		errorHandler.HandleCatalogError(ctx, c, err, slog.Int("category_id", id))
		return
	}

	slog.InfoContext(ctx, "category found",
		slog.Group("category",
			slog.Int("id", category.ID()),
			slog.String("name", category.Name()),
		),
	)
	c.JSON(http.StatusOK, generated.GetCategoryById200JSONResponse{
		Data: mapper.ToGeneratedCategoryResponse(category),
	})
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context, id generated.ID) {
	ctx := c.Request.Context()
	if !validation.ValidID(id) {
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: entities.ErrInvalidCategoryID.Error(),
		})
		return
	}

	var request generated.UpdateCategoryJSONRequestBody
	if err := c.ShouldBindJSON(&request); err != nil {
		errorHandler.HandleCatalogError(ctx, c, err, slog.Int("category_id", id))
		return
	}

	categoryCommand := command.UpdateCategoryCommand{
		ID:   id,
		Name: request.Name,
	}

	category, err := h.service.UpdateCategory(ctx, &categoryCommand)
	if err != nil {
		attrs := []slog.Attr{slog.Int("category_id", categoryCommand.ID)}
		if request.Name != nil {
			attrs = append(attrs, slog.String("category_name", *categoryCommand.Name))
		}
		errorHandler.HandleCatalogError(ctx, c, err, attrs...)
		return
	}

	slog.InfoContext(ctx, "updated category",
		slog.Group("category",
			slog.Int("id", id),
			slog.String("name", category.Name()),
		),
	)
	c.JSON(http.StatusOK, generated.UpdateCategory200JSONResponse{
		Data: mapper.ToGeneratedCategoryResponse(category),
	})
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context, id generated.ID) {
	ctx := c.Request.Context()
	if !validation.ValidID(id) {
		c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: entities.ErrInvalidCategoryID.Error(),
		})
		return
	}

	err := h.service.DeleteCategory(ctx, id)
	if err != nil {
		errorHandler.HandleCatalogError(ctx, c, err, slog.Int("category_id", id))
		return
	}

	slog.InfoContext(ctx, "category deleted",
		slog.Int("category_id", id),
	)
	c.Status(http.StatusNoContent)
}
