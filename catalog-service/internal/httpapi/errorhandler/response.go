package errorhandler

import (
	"context"
	"log/slog"
	"virtual-library/catalog-service/internal/httpapi/generated"

	"github.com/gin-gonic/gin"
)

func writeError(
	ctx context.Context,
	c *gin.Context,
	status int,
	level slog.Level,
	message string,
	attrs ...slog.Attr,
) {
	slog.LogAttrs(ctx, level, message, attrs...)

	c.JSON(status, generated.ErrorResponse{
		Message: message,
	})
}
