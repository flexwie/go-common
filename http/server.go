package http

import (
	"context"
	"net"
	"net/http"

	"go.uber.org/fx"
)

var WithHttp = fx.Provide(
	newHttpServer,
	fx.Annotate(
		newServeMux,
		fx.ParamTags(`group:"routes"`),
	),
)

func newHttpServer(lc fx.Lifecycle, mux *http.ServeMux) *http.Server {
	s := &http.Server{Addr: ":8080", Handler: mux}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", s.Addr)
			if err != nil {
				return err
			}

			go s.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return s.Shutdown(ctx)
		},
	})

	return s
}

func newServeMux(routes []Route) *http.ServeMux {
	mux := http.NewServeMux()
	for _, route := range routes {
		mux.Handle(route.Pattern(), route)
	}
	return mux
}
