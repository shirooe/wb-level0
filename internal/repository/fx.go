package repository

import (
	"context"
	"wb-level0/internal/database"
	"wb-level0/internal/models"

	"go.uber.org/fx"
)

type Repository interface {
	GetAllOrders(ctx context.Context) ([]models.Order, error)

	GetOrderByID(ctx context.Context, id string) (models.Order, error)
	GetDeliveryByID(ctx context.Context, id string) (models.Delivery, error)
	GetPaymentByID(ctx context.Context, id string) (models.Payment, error)
	GetItemsByID(ctx context.Context, id string) ([]models.Item, error)

	CreateOrder(ctx context.Context, order models.Order) (string, error)
	CreateDelivery(ctx context.Context, orderID string, delivery models.Delivery) error
	CreatePayment(ctx context.Context, orderID string, payment models.Payment) error
	CreateItems(ctx context.Context, orderID string, items []models.Item) error
}

type repository struct {
	db *database.Client
}

func NewRepository(db *database.Client) Repository {
	return &repository{
		db: db,
	}
}

// регистрация модуля
func Module() fx.Option {
	return fx.Module("repository",
		fx.Provide(NewRepository))
}
