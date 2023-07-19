package grpc

import (
	"context"
	"time"

	"github.com/charmbracelet/log"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

var WithUnaryLoggingInterceptorFactory = fx.Options(
	fx.Provide(fx.Annotate(
		NewUnaryLoggingInterceptorFactory,
		fx.ResultTags("unary_logging"),
	)),
)

func NewUnaryLoggingInterceptorFactory(logger *log.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		h, err := handler(ctx, req)

		if err != nil {
			logger.Warn(info.FullMethod, "time", time.Since(start), "err", err)
		} else {
			logger.Info(info.FullMethod, "time", time.Since(start), "err", err)
		}

		return h, err
	}
}
