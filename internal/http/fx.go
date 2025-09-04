package http

import (
	"context"
	"net"
	"net/http"
	"wb-level0/internal/service"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// регистрация модуля http
func Module() fx.Option {
	return fx.Module("http",
		// регистрация конфига, роутера и контроллера
		fx.Provide(ProvideConfig, NewServerMux, NewController),
		// запуск модуля
		fx.Invoke(func(lc fx.Lifecycle, config *Config, mux *mux.Router, controller *Controller, service *service.WBLevel0Service, log *zap.Logger) {
			// регистрация роутов
			controller.RegisterRoutes(mux)

			// создание экземпляра сервера
			server := &http.Server{
				// адрес сервера
				Addr: config.Address(),
				// настройка CORS
				Handler: handlers.CORS(
					handlers.AllowedOrigins([]string{"*"}),
					handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
					handlers.AllowedHeaders([]string{"*"}),
				)(mux),
			}

			// жизненный цикл модуля
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					ln, err := net.Listen("tcp", server.Addr)
					if err != nil {
						return err
					}

					go func() {
						// запуск сервера
						if err := server.Serve(ln); err != nil {
							// TODO: Error starting or closing listener
						}
					}()

					// восстановление заказов в кэш из БД
					if err := service.RestoreOrders(ctx); err != nil {
						return err
					}

					log.Info("[http] запуск сервера")
					return nil
				},
				OnStop: func(ctx context.Context) error {
					log.Info("[http] остановка сервера")
					// остановка сервера
					return server.Shutdown(ctx)
				},
			})
		}))
}
