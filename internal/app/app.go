package app

import (
	"context"
	"wb-level0/internal/cache"
	"wb-level0/internal/config"
	"wb-level0/internal/database"
	"wb-level0/internal/http"
	"wb-level0/internal/kafka"
	"wb-level0/internal/repository"
	"wb-level0/internal/service"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func New() *fx.App {
	return fx.New(
		fx.NopLogger,
		fx.Provide(func() context.Context {
			return context.Background()
		}),
		fx.Provide(config.New, zap.NewDevelopment),
		fx.Options(http.Module(), kafka.Module(),
			service.Module(), repository.Module(), database.Module(), cache.Module()),
	)
}
