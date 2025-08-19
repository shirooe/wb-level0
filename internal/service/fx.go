package service

import (
	"wb-level0/internal/repository"

	"go.uber.org/fx"
)

type WBLevel0Service struct {
	repository repository.Repository
}

func NewWBLevel0Service(repository repository.Repository) *WBLevel0Service {
	return &WBLevel0Service{
		repository: repository,
	}
}

func Module() fx.Option {
	return fx.Module("service",
		fx.Provide(NewWBLevel0Service))
}
