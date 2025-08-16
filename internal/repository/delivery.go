package repository

import (
	"context"
	"fmt"
	"wb-level0/internal/database"
	"wb-level0/internal/models"

	"github.com/elgris/sqrl"
)

type Delivery struct {
	db *database.Client
}

func NewDeliveryRepository(db *database.Client) *Delivery {
	return &Delivery{
		db: db,
	}
}

func (o *Delivery) Create(ctx context.Context, orderID string, delivery models.Delivery) (string, error) {
	sql, args, err := sqrl.Insert("delivery").PlaceholderFormat(sqrl.Dollar).Columns("order_uid", "name", "phone", "zip", "city", "address", "region", "email").
		Values(orderID, delivery.Name, delivery.Phone, delivery.Zip, delivery.City, delivery.Address, delivery.Region, delivery.Email).
		ToSql()

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	query := database.Query{
		Name:     "CreateDelivery",
		QueryRaw: sql,
	}

	o.db.DB().QueryRowContext(ctx, query, args...)
	return orderID, nil
}
