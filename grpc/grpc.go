package grpc

import (
	"context"
	"net"
	"net/url"
	"time"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

var WithGrpc = fx.Provide(newGrpc)

func newGrpc(lc fx.Lifecycle) *grpc.Server {
	addr := "tcp://0.0.0.0:9000"
	uri, err := url.Parse(addr)
	if err != nil {
		panic(err)
	}

	ln, err := net.Listen(uri.Scheme, uri.Host)
	if err != nil {
		panic(err)
	}

	keepaliveOpts := grpc.KeepaliveParams(keepalive.ServerParameters{
		Time:    time.Minute,
		Timeout: 3 * time.Second,
	})

	keepaliveEnforcement := grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
		MinTime:             30 * time.Second,
		PermitWithoutStream: true,
	})

	server := grpc.NewServer(keepaliveOpts, keepaliveEnforcement)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				reflection.Register(server)

				if err := server.Serve(ln); err != nil {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			server.GracefulStop()
			return nil
		},
	})

	return server
}
