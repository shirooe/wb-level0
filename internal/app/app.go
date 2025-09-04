package app

import (
	"context"
	"wb-level0/internal/cache"
	"wb-level0/internal/config"
	"wb-level0/internal/database"
	"wb-level0/internal/database/transaction"
	"wb-level0/internal/http"
	"wb-level0/internal/kafka"
	"wb-level0/internal/repository"
	"wb-level0/internal/service"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// регистрация необходимых модулей
func New() *fx.App {
	return fx.New(
		// отключение логгера fx
		fx.NopLogger,
		// регистрация контекст для других приложении
		fx.Provide(func() context.Context {
			return context.Background()
		}),
		// регистрация конфига и логгера zap
		fx.Provide(config.New, zap.NewDevelopment),
		fx.Options(database.Module(), http.Module(), kafka.Module(),
			service.Module(), repository.Module(), transaction.Module(), cache.Module()),
	)
}
