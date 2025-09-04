package transaction

import (
	"wb-level0/internal/database"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Manager struct {
	db  *database.Client
	log *zap.Logger
}

// создание менеджера транзакции
func NewManager(db *database.Client, log *zap.Logger) *Manager {
	return &Manager{
		db:  db,
		log: log,
	}
}

func Module() fx.Option {
	return fx.Module("transaction",
		fx.Provide(NewManager),
	)
}
