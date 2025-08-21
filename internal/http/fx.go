package http

import (
	"context"
	"net"
	"net/http"
	"wb-level0/internal/service"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Module() fx.Option {
	return fx.Module("http",
		fx.Provide(ProvideConfig, NewServerMux, NewController),
		fx.Invoke(func(lc fx.Lifecycle, config *Config, mux *mux.Router, controller *Controller, service *service.WBLevel0Service, log *zap.Logger) {
			controller.RegisterRoutes(mux)

			server := &http.Server{
				Addr:    config.Address(),
				Handler: mux,
			}

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					ln, err := net.Listen("tcp", server.Addr)
					if err != nil {
						return err
					}

					go func() {
						if err := server.Serve(ln); err != nil {
							// TODO: Error starting or closing listener
						}
					}()

					if err := service.RestoreOrders(ctx); err != nil {
						return err
					}

					log.Info("[http] запуск сервера")
					return nil
				},
				OnStop: func(ctx context.Context) error {
					log.Info("[http] остановка сервера")
					return server.Shutdown(ctx)
				},
			})
		}))
}
