package middleware

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

func RequestLogAttrs(c *gin.Context, err error) []slog.Attr {
	return []slog.Attr{
		slog.String("method", c.Request.Method),
		slog.String("path", c.Request.URL.Path),
		slog.String("error", err.Error()),
	}
}
