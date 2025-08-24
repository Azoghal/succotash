package log

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func ZerologMiddleware(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next() // Process request

		// Log request details
		logger.Info().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Dur("duration", time.Since(start)).
			Msg("Request completed")
	}
}