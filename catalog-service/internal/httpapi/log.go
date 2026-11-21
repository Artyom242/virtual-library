package httpapi

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

func requestLogAttrs(c *gin.Context) []slog.Attr {
	return []slog.Attr{
		slog.String("method", c.Request.Method),
		slog.String("path", c.Request.URL.Path),
	}
}
