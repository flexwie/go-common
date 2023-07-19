package logger

import (
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
)

func NewGinLogger(logger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now() // Start timer
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Fill the params
		param := gin.LogFormatterParams{}

		param.TimeStamp = time.Now() // Stop timer
		param.Latency = param.TimeStamp.Sub(start)
		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}

		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		param.BodySize = c.Writer.Size()
		if raw != "" {
			path = path + "?" + raw
		}
		param.Path = path

		poplogger := logger.With("time", param.TimeStamp, "method", param.Method, "status_code", param.StatusCode, "body_size", param.BodySize, "duration", param.Latency).WithPrefix("gin")

		if c.Writer.Status() >= 500 {
			poplogger.Error(param.Path)
		} else {
			poplogger.Info(param.Path)
		}
	}
}
