package http

import (
	"context"
	"net"
	"net/http"

	"github.com/charmbracelet/log"
	"go.uber.org/fx"
)

func WithHttpFactory(addr string) fx.Option {
	return fx.Provide(
		newHttpServerFactory(addr),
		fx.Annotate(
			newServeMux,
			fx.ParamTags(`group:"routes"`),
		),
	)
}

var WithHttp = fx.Provide(
	newHttpServer,
	fx.Annotate(
		newServeMux,
		fx.ParamTags(`group:"routes"`),
	),
)

func newHttpServerFactory(addr string) func(fx.Lifecycle, *http.ServeMux, *log.Logger) *http.Server {
	return func(lc fx.Lifecycle, mux *http.ServeMux, logger *log.Logger) *http.Server {
		s := &http.Server{Addr: addr, Handler: mux}

		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				ln, err := net.Listen("tcp", s.Addr)
				if err != nil {
					return err
				}

				go s.Serve(ln)
				logger.Info("started http server", "addr", s.Addr)

				return nil
			},
			OnStop: func(ctx context.Context) error {
				return s.Shutdown(ctx)
			},
		})

		return s
	}
}

func newHttpServer(lc fx.Lifecycle, mux *http.ServeMux, logger *log.Logger) *http.Server {
	s := &http.Server{Addr: ":8080", Handler: mux}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", s.Addr)
			if err != nil {
				return err
			}

			go s.Serve(ln)
			logger.Info("started http server", "addr", s.Addr)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return s.Shutdown(ctx)
		},
	})

	return s
}

func newServeMux(routes []Route, logger *log.Logger) *http.ServeMux {
	logger = logger.WithPrefix("routing")

	mux := http.NewServeMux()
	for _, route := range routes {
		logger.Debug("adding route", "from", route.Pattern(), "to", route)
		mux.Handle(route.Pattern(), route)
	}
	return mux
}
