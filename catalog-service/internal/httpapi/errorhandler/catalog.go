package errorhandler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"virtual-library/catalog-service/internal/app"
	"virtual-library/catalog-service/internal/domain/entity"
	"virtual-library/catalog-service/internal/httpapi/middleware"

	"github.com/gin-gonic/gin"
)

func HandleCatalogError(
	ctx context.Context,
	c *gin.Context,
	err error,
	attrs ...slog.Attr,
) {
	attrs = append(middleware.RequestLogAttrs(c, err), attrs...)

	var (
		status  int
		level   slog.Level
		message string
	)

	switch {
	case errors.Is(err, io.EOF):
		status = http.StatusBadRequest
		level = slog.LevelWarn
		message = ErrInvalidRequest

	case isJSONError(err):
		status = http.StatusBadRequest
		level = slog.LevelWarn
		message = ErrInvalidRequest

	case errors.Is(err, entity.ErrValidationError):
		status = http.StatusBadRequest
		level = slog.LevelWarn
		message = entity.ErrValidationError.Error()

	// Book errors
	case errors.Is(err, entity.ErrBookNotFound):
		status = http.StatusNotFound
		level = slog.LevelWarn
		message = entity.ErrBookNotFound.Error()

	case errors.Is(err, entity.ErrInvalidBookData):
		status = http.StatusBadRequest
		level = slog.LevelWarn
		message = entity.ErrInvalidBookData.Error()

	case errors.Is(err, entity.ErrInvalidBookDescription):
		status = http.StatusBadRequest
		level = slog.LevelWarn
		message = entity.ErrInvalidBookDescription.Error()

	case errors.Is(err, app.ErrEmptyBookPatch):
		status = http.StatusBadRequest
		level = slog.LevelWarn
		message = app.ErrEmptyBookPatch.Error()

	// Author errors
	case errors.Is(err, entity.ErrAuthorNotFound):
		status = http.StatusNotFound
		level = slog.LevelWarn
		message = entity.ErrAuthorNotFound.Error()

	case errors.Is(err, entity.ErrInvalidAuthorData):
		status = http.StatusBadRequest
		level = slog.LevelWarn
		message = entity.ErrInvalidAuthorData.Error()

	case errors.Is(err, app.ErrAuthorAlreadyExists):
		status = http.StatusConflict
		level = slog.LevelWarn
		message = app.ErrAuthorAlreadyExists.Error()

	case errors.Is(err, app.ErrEmptyAuthorPatch):
		status = http.StatusBadRequest
		level = slog.LevelWarn
		message = app.ErrEmptyAuthorPatch.Error()

	case errors.Is(err, app.ErrAuthorHasBooks):
		status = http.StatusConflict
		level = slog.LevelWarn
		message = app.ErrAuthorHasBooks.Error()

	// Category errors
	case errors.Is(err, entity.ErrCategoryNotFound):
		status = http.StatusNotFound
		level = slog.LevelWarn
		message = entity.ErrCategoryNotFound.Error()

	case errors.Is(err, entity.ErrCategoryAlreadyExists):
		status = http.StatusConflict
		level = slog.LevelWarn
		message = entity.ErrCategoryAlreadyExists.Error()

	default:
		status = http.StatusInternalServerError
		level = slog.LevelError
		message = ErrInternalServerMessage

		attrs = append(attrs, slog.String("error", err.Error()))
	}

	writeError(ctx, c, status, level, message, attrs...)
}

func isJSONError(err error) bool {
	var syntaxErr *json.SyntaxError
	var typeErr *json.UnmarshalTypeError

	return errors.As(err, &syntaxErr) ||
		errors.As(err, &typeErr)
}
