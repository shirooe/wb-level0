package repository

import (
	"context"
	"wb-level0/internal/database"
	"wb-level0/internal/models"

	"github.com/elgris/sqrl"
)

// создание предметов в БД
func (r *repository) CreateItems(ctx context.Context, orderID string, items []models.Item) error {
	// цикл для каждого предмета
	for _, item := range items {
		// создание sql строки и аргументов к нему
		sql, args, err := sqrl.Insert("items").PlaceholderFormat(sqrl.Dollar).Columns("order_uid", "chrt_id", "track_number", "price", "rid", "name", "sale", "size",
			"total_price", "nm_id", "brand", "status").
			Values(orderID, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size,
				item.TotalPrice, item.NmID, item.Brand, item.Status).
			ToSql()

		if err != nil {
			return err
		}

		// создание структуры запроса
		query := database.Query{
			Name:     "CreateItems",
			QueryRaw: sql,
		}

		// выполнение запроса
		r.db.DB().QueryRowContext(ctx, query, args...)
	}

	return nil
}

// поиск предмета по схождению order_uid
func (r *repository) GetItemsByID(ctx context.Context, id string) ([]models.Item, error) {
	// создание sql строки и аргументов к нему
	sql, args, err := sqrl.Select("order_uid", "chrt_id", "track_number", "price", "rid", "name", "sale", "size",
		"total_price", "nm_id", "brand", "status").
		PlaceholderFormat(sqrl.Dollar).From("items").Where(sqrl.Eq{"order_uid": id}).ToSql()

	if err != nil {
		return nil, err
	}

	// создание структуры запроса
	query := database.Query{
		Name:     "GetItems",
		QueryRaw: sql,
	}

	// сканирование всех предметов в слайс структуру
	var items []models.Item
	if err := r.db.DB().ScanAllContext(ctx, &items, query, args...); err != nil {
		return nil, err
	}

	// возвращение данных
	return items, nil
}
