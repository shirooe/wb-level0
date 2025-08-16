package repository

import (
	"context"
	"wb-level0/internal/database"
	"wb-level0/internal/models"

	"github.com/elgris/sqrl"
)

type Item struct {
	db *database.Client
}

func NewItemRepository(db *database.Client) *Item {
	return &Item{
		db: db,
	}
}

func (o *Item) Create(ctx context.Context, orderID string, items []models.Item) (string, error) {
	for _, item := range items {
		sql, args, err := sqrl.Insert("items").PlaceholderFormat(sqrl.Dollar).Columns("order_uid", "chrt_id", "track_number", "price", "rid", "name", "sale", "size",
			"total_price", "nm_id", "brand", "status").
			Values(orderID, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size,
				item.TotalPrice, item.NmID, item.Brand, item.Status).
			ToSql()

		if err != nil {
			return "", err
		}

		query := database.Query{
			Name:     "CreateItem",
			QueryRaw: sql,
		}

		o.db.DB().QueryRowContext(ctx, query, args...)
	}

	return orderID, nil
}
