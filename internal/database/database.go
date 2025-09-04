package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// клиент для взаимодействия в БД
type Client struct {
	db     DB
	config *Config
}

// подключение к БД 
func ProvideClient(ctx context.Context, cfg *Config) (*Client, error) {
	conn, err := pgxpool.New(ctx, cfg.DSN())

	if err != nil {
		return nil, err
	}

	return &Client{
		db:     pg(conn),
		config: cfg,
	}, nil
}

// композиция интерфейсов для взаимойдействия с БД
func (c *Client) DB() DB {
	return c.db
}

// Пинг БД
func (c *Client) Ping(ctx context.Context) error {
	return c.db.Ping(ctx)
}

// Закрытие соединения
func (c *Client) Close() error {
	return c.db.Close()
}
