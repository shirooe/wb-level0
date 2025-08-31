package service

import (
	"wb-level0/internal/cache"
	"wb-level0/internal/database/transaction"
	"wb-level0/internal/repository"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type WBLevel0Service struct {
	repository repository.Repository
	manager    *transaction.Manager
	cache      *cache.Cache
	log        *zap.Logger
}

func NewWBLevel0Service(repository repository.Repository, manager *transaction.Manager, cache *cache.Cache, log *zap.Logger) *WBLevel0Service {
	return &WBLevel0Service{
		repository: repository,
		manager:    manager,
		cache:      cache,
		log:        log,
	}
}

func Module() fx.Option {
	return fx.Module("service",
		fx.Provide(NewWBLevel0Service))
}
