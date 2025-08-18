package repository

import (
	"context"
	"wb-level0/internal/database"
	"wb-level0/internal/models"

	"go.uber.org/fx"
)

type Repository interface {
	CreateOrder(ctx context.Context, order models.Order) (string, error)
	CreateItem(ctx context.Context, orderID string, items []models.Item) (string, error)
	CreatePayment(ctx context.Context, orderID string, payment models.Payment) (string, error)
	CreateDelivery(ctx context.Context, orderID string, delivery models.Delivery) (string, error)
}

type repository struct {
	db *database.Client
}

func NewRepository(db *database.Client) Repository {
	return &repository{
		db: db,
	}
}

func Module() fx.Option {
	return fx.Module("repository",
		fx.Provide(NewRepository))
}
