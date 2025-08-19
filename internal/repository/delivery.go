package repository

import (
	"context"
	"wb-level0/internal/database"
	"wb-level0/internal/models"

	"github.com/elgris/sqrl"
)

func (r *repository) CreateDelivery(ctx context.Context, orderID string, delivery models.Delivery) error {
	sql, args, err := sqrl.Insert("delivery").PlaceholderFormat(sqrl.Dollar).Columns("order_uid", "name", "phone", "zip", "city", "address", "region", "email").
		Values(orderID, delivery.Name, delivery.Phone, delivery.Zip, delivery.City, delivery.Address, delivery.Region, delivery.Email).
		ToSql()

	if err != nil {
		return err
	}

	query := database.Query{
		Name:     "CreateDelivery",
		QueryRaw: sql,
	}

	r.db.DB().QueryRowContext(ctx, query, args...)
	return nil
}

func (r *repository) GetDeliveryByID(ctx context.Context, id string) (models.Delivery, error) {
	sql, args, err := sqrl.Select("order_uid", "name", "phone", "zip", "city", "address", "region", "email").
		PlaceholderFormat(sqrl.Dollar).From("delivery").Where(sqrl.Eq{"order_uid": id}).ToSql()

	if err != nil {
		return models.Delivery{}, err
	}

	query := database.Query{
		Name:     "GetDelivery",
		QueryRaw: sql,
	}

	var delivery models.Delivery
	if err := r.db.DB().ScanOneContext(ctx, &delivery, query, args...); err != nil {
		return models.Delivery{}, err
	}

	return delivery, nil
}
