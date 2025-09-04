package repository

import (
	"context"
	"wb-level0/internal/database"
	"wb-level0/internal/models"

	"github.com/elgris/sqrl"
)

// создание доставки в БД
func (r *repository) CreateDelivery(ctx context.Context, orderID string, delivery models.Delivery) error {
	// создание sql строки и аргументов к нему
	sql, args, err := sqrl.Insert("delivery").PlaceholderFormat(sqrl.Dollar).Columns("order_uid", "name", "phone", "zip", "city", "address", "region", "email").
		Values(orderID, delivery.Name, delivery.Phone, delivery.Zip, delivery.City, delivery.Address, delivery.Region, delivery.Email).
		ToSql()

	if err != nil {
		return err
	}

	// создание структуры запроса
	query := database.Query{
		Name:     "CreateDelivery",
		QueryRaw: sql,
	}

	// выполнение запроса
	r.db.DB().QueryRowContext(ctx, query, args...)
	return nil
}

// поиск доставки по order_uid в БД
func (r *repository) GetDeliveryByID(ctx context.Context, id string) (models.Delivery, error) {
	// создание sql строки и аргументов к нему
	sql, args, err := sqrl.Select("order_uid", "name", "phone", "zip", "city", "address", "region", "email").
		PlaceholderFormat(sqrl.Dollar).From("delivery").Where(sqrl.Eq{"order_uid": id}).ToSql()

	if err != nil {
		return models.Delivery{}, err
	}

	// создание структуры запроса
	query := database.Query{
		Name:     "GetDelivery",
		QueryRaw: sql,
	}

	// сканирование в структуру
	var delivery models.Delivery
	if err := r.db.DB().ScanOneContext(ctx, &delivery, query, args...); err != nil {
		return models.Delivery{}, err
	}

	// возвращение данных
	return delivery, nil
}
