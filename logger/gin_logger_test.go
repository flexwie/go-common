package logger_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/charmbracelet/log"
	clogger "github.com/flexwie/go-common/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLogRoute(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf)
	w := httptest.NewRecorder()

	r := gin.New()
	gin.SetMode(gin.ReleaseMode)
	r.Use(clogger.NewGinLogger(logger))
	r.Handle("GET", "/", func(ctx *gin.Context) {
	})

	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, buf.String(), "INFO")
	assert.Contains(t, buf.String(), "gin: /")
	assert.Contains(t, buf.String(), "method=GET")
}

func TestLogQuery(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf)
	w := httptest.NewRecorder()

	r := gin.New()
	gin.SetMode(gin.ReleaseMode)
	r.Use(clogger.NewGinLogger(logger))
	r.Handle("GET", "/", func(ctx *gin.Context) {
	})

	req, _ := http.NewRequest("GET", "/?test=test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, buf.String(), "INFO")
	assert.Contains(t, buf.String(), "gin: /?test=test")
	assert.Contains(t, buf.String(), "method=GET")
}

func TestLogError(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf)
	w := httptest.NewRecorder()

	r := gin.New()
	gin.SetMode(gin.ReleaseMode)
	r.Use(clogger.NewGinLogger(logger))
	r.Handle("GET", "/", func(ctx *gin.Context) {
		ctx.AbortWithError(502, errors.New("fail"))
	})

	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 502, w.Code)
	assert.Contains(t, buf.String(), "ERRO")
	assert.Contains(t, buf.String(), "gin: /")
	assert.Contains(t, buf.String(), "method=GET")
}
