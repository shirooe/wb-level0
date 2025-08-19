package repository

import (
	"context"
	"wb-level0/internal/database"
	"wb-level0/internal/models"

	"github.com/elgris/sqrl"
)

func (r *repository) CreateOrder(ctx context.Context, order models.Order) (string, error) {
	sql, args, err := sqrl.Insert("orders").PlaceholderFormat(sqrl.Dollar).Columns("order_uid", "track_number", "entry",
		"locale", "internal_signature", "customer_id", "delivery_service", "shardkey", "sm_id", "date_created", "oof_shard",
	).Values(order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID,
		order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OofShard).Returning("order_uid").ToSql()

	if err != nil {
		return "", err
	}

	query := database.Query{
		Name:     "CreateOrder",
		QueryRaw: sql,
	}

	var id string
	if err := r.db.DB().QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}

func (r *repository) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	var orders []models.Order

	sql, args, err := sqrl.Select("order_uid").PlaceholderFormat(sqrl.Dollar).From("orders").ToSql()
	if err != nil {
		return nil, err
	}

	query := database.Query{
		Name:     "GetAllOrders",
		QueryRaw: sql,
	}

	rows, err := r.db.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order_uid string
		if err := rows.Scan(&order_uid); err != nil {
			return nil, err
		}

		order, err := r.GetOrderByID(ctx, order_uid)
		if err != nil {
			return nil, err
		}

		delivery, err := r.GetDeliveryByID(ctx, order_uid)
		if err != nil {
			return nil, err
		}

		payment, err := r.GetPaymentByID(ctx, order_uid)
		if err != nil {
			return nil, err
		}

		items, err := r.GetItemsByID(ctx, order_uid)
		if err != nil {
			return nil, err
		}

		order.Delivery = delivery
		order.Payment = payment
		order.Items = items

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *repository) GetOrderByID(ctx context.Context, id string) (models.Order, error) {
	sql, args, err := sqrl.Select("order_uid", "track_number", "entry", "locale", "internal_signature", "customer_id",
		"delivery_service", "shardkey", "sm_id", "date_created", "oof_shard").
		PlaceholderFormat(sqrl.Dollar).From("orders").Where(sqrl.Eq{"order_uid": id}).ToSql()

	if err != nil {
		return models.Order{}, err
	}

	query := database.Query{
		Name:     "GetOrder",
		QueryRaw: sql,
	}

	var order models.Order
	if err := r.db.DB().ScanOneContext(ctx, &order, query, args...); err != nil {
		return models.Order{}, err
	}

	return order, nil
}
