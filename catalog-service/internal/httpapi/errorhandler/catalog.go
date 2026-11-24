package errorhandler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"virtual-library/catalog-service/internal/app"
	"virtual-library/catalog-service/internal/domain/entities"
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

	case errors.Is(err, entities.ErrValidationError):
		status = http.StatusBadRequest
		level = slog.LevelWarn
		message = entities.ErrValidationError.Error()

	// Book errors
	case errors.Is(err, entities.ErrBookNotFound):
		status = http.StatusNotFound
		level = slog.LevelWarn
		message = entities.ErrBookNotFound.Error()

	case errors.Is(err, entities.ErrInvalidBookData):
		status = http.StatusBadRequest
		level = slog.LevelWarn
		message = entities.ErrInvalidBookData.Error()

	case errors.Is(err, entities.ErrInvalidBookDescription):
		status = http.StatusBadRequest
		level = slog.LevelWarn
		message = entities.ErrInvalidBookDescription.Error()

	case errors.Is(err, app.ErrEmptyBookPatch):
		status = http.StatusBadRequest
		level = slog.LevelWarn
		message = app.ErrEmptyBookPatch.Error()

	// Author errors
	case errors.Is(err, entities.ErrAuthorNotFound):
		status = http.StatusNotFound
		level = slog.LevelWarn
		message = entities.ErrAuthorNotFound.Error()

	case errors.Is(err, entities.ErrInvalidAuthorData):
		status = http.StatusBadRequest
		level = slog.LevelWarn
		message = entities.ErrInvalidAuthorData.Error()

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
	case errors.Is(err, entities.ErrCategoryNotFound):
		status = http.StatusNotFound
		level = slog.LevelWarn
		message = entities.ErrCategoryNotFound.Error()

	case errors.Is(err, entities.ErrCategoryAlreadyExists):
		status = http.StatusConflict
		level = slog.LevelWarn
		message = entities.ErrCategoryAlreadyExists.Error()

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
