package repository

import (
	"context"
	"wb-level0/internal/database"
	"wb-level0/internal/models"

	"github.com/elgris/sqrl"
)

// создание заказа в БД
func (r *repository) CreateOrder(ctx context.Context, order models.Order) (string, error) {
	// создание sql строки и аргументов к нему
	sql, args, err := sqrl.Insert("orders").PlaceholderFormat(sqrl.Dollar).Columns("order_uid", "track_number", "entry",
		"locale", "internal_signature", "customer_id", "delivery_service", "shardkey", "sm_id", "date_created", "oof_shard",
	).Values(order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID,
		order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OofShard).Returning("order_uid").ToSql()

	if err != nil {
		return "", err
	}

	// создание структуры запроса
	query := database.Query{
		Name:     "CreateOrder",
		QueryRaw: sql,
	}

	// сканирование order_uid в переменную
	var id string
	if err := r.db.DB().ScanOneContext(ctx, &id, query, args...); err != nil {
		return "", err
	}

	// возвращение order_uid
	return id, nil
}

// поиск всех заказов в БД для кэша
func (r *repository) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	var orders []models.Order

	// создание sql строки и аргументов к нему
	sql, args, err := sqrl.Select("order_uid").PlaceholderFormat(sqrl.Dollar).From("orders").ToSql()
	if err != nil {
		return nil, err
	}

	// создание структуры запроса
	query := database.Query{
		Name:     "GetAllOrders",
		QueryRaw: sql,
	}

	// выполнение запроса и получение всех строк
	rows, err := r.db.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	// закрытие
	defer rows.Close()

	// FIX ME
	for rows.Next() {
		// сканирование каждой строки
		var order_uid string
		if err := rows.Scan(&order_uid); err != nil {
			continue
		}

		// получение заказа
		order, err := r.GetOrderByID(ctx, order_uid)
		if err != nil {
			continue
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// возвращение данных
	return orders, nil
}

// поиск заказа в БД по order_uid
func (r *repository) GetOrderByID(ctx context.Context, id string) (models.Order, error) {
	// создание sql строки и аргументов к нему
	sql, args, err := sqrl.Select("order_uid", "track_number", "entry", "locale", "internal_signature", "customer_id",
		"delivery_service", "shardkey", "sm_id", "date_created", "oof_shard").
		PlaceholderFormat(sqrl.Dollar).From("orders").Where(sqrl.Eq{"order_uid": id}).ToSql()

	if err != nil {
		return models.Order{}, err
	}

	// создание структуры запроса
	query := database.Query{
		Name:     "GetOrder",
		QueryRaw: sql,
	}

	// получение одного заказа
	var order models.Order
	if err := r.db.DB().ScanOneContext(ctx, &order, query, args...); err != nil {
		return models.Order{}, err
	}

	// поиск доставки
	delivery, err := r.GetDeliveryByID(ctx, id)
	if err != nil {
		return models.Order{}, err
	}

	// поиск оплаты
	payment, err := r.GetPaymentByID(ctx, id)
	if err != nil {
		return models.Order{}, err
	}

	// поиск предметов
	items, err := r.GetItemsByID(ctx, id)
	if err != nil {
		return models.Order{}, err
	}

	order.Delivery = delivery
	order.Payment = payment
	order.Items = items

	// возвращение данных
	return order, nil
}
