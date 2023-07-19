package grpc_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/charmbracelet/log"
	cgrpc "github.com/flexwie/go-common/grpc"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type empty struct{}

func TestCanLogError(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf)

	i := cgrpc.NewUnaryLoggingInterceptorFactory(logger)

	i(context.TODO(), empty{}, &grpc.UnaryServerInfo{}, func(c context.Context, req interface{}) (interface{}, error) {
		return empty{}, errors.New("sample")
	})

	assert.Contains(t, buf.String(), "err=sample")
	assert.Contains(t, buf.String(), "WARN")
}

func TestCanLogName(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf)

	i := cgrpc.NewUnaryLoggingInterceptorFactory(logger)

	i(context.TODO(), empty{}, &grpc.UnaryServerInfo{FullMethod: "testMethod"}, func(c context.Context, req interface{}) (interface{}, error) {
		return empty{}, nil
	})

	assert.Contains(t, buf.String(), "testMethod")
	assert.Contains(t, buf.String(), "INFO")
}
