package httpapi

import (
	"errors"
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
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
		})
		return
	}

	author := domain.Author{
		Name: request.Name,
	}

	author, err := h.service.Create(ctx, author)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "created",
		"data":    author,
	})
}
func (h *AuthorHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	authors, err := h.service.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": authors,
	})
}
func (h *AuthorHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	request := c.Param("id")
	id, err := strconv.Atoi(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
		})
		return
	}
	author, err := h.service.GetById(ctx, id)

	if errors.Is(err, domain.ErrBookNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Book not found",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": author,
	})
}
func (h *AuthorHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()
	data := c.Param("id")
	id, err := strconv.Atoi(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
		})
		return
	}

	var request dto.UpdateAuthorRequest
	if err = c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
		})
		return
	}

	patch := domain.AuthorPatch{Name: request.Name}
	author, err := h.service.Update(ctx, id, patch)
	if errors.Is(err, domain.ErrAuthorNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Author not found",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": author,
	})
}
func (h *AuthorHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
		})
		return
	}

	err = h.service.Delete(ctx, id)
	if errors.Is(err, domain.ErrBookNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Book not found",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{
		"message": "deleted",
	})
}
