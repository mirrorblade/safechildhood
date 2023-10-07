package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (h *Handler) loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		codeStatus := c.Writer.Status()

		end := time.Now()

		fields := []zapcore.Field{
			zap.Int("status", codeStatus),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Duration("latency", time.Since(start)),
			zap.String("time", end.Format(time.RFC3339)),
		}

		if len(c.Errors) > 0 {
			for _, err := range c.Errors.Errors() {
				h.logger.Error(err, fields...)
			}

		} else {
			if codeStatus < 400 {
				h.logger.Info(path, fields...)

			} else if codeStatus < 500 {
				h.logger.Warn(path, fields...)

			} else {
				h.logger.Error(path, fields...)
			}
		}
	}
}
