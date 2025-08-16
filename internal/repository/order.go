package repository

import (
	"context"
	"fmt"
	"wb-level0/internal/database"
	"wb-level0/internal/models"

	"github.com/elgris/sqrl"
)

type Order struct {
	db *database.Client
}

func NewOrderRepository(db *database.Client) *Order {
	return &Order{
		db: db,
	}
}

func (o *Order) Create(ctx context.Context, order models.Order) (string, error) {
	sql, args, err := sqrl.Insert("orders").PlaceholderFormat(sqrl.Dollar).Columns("order_uid", "track_number", "entry",
		"locale", "internal_signature", "customer_id", "delivery_service", "shardkey", "sm_id", "date_created", "oof_shard",
	).Values(order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID,
		order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OofShard).Returning("order_uid").ToSql()

	if err != nil {
		return "", err
	}

	fmt.Printf("arg %v\n", args)

	query := database.Query{
		Name:     "CreateOrder",
		QueryRaw: sql,
	}

	var id string
	if err := o.db.DB().QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		fmt.Printf("Error creating order: %v\n", err)
		return "", err
	}

	fmt.Printf("Order created with ID: %s\n", id)
	return id, nil
}
