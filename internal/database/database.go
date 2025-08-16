package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	db     DB
	config Config
}

func New(ctx context.Context, cfg *Config) (*Client, error) {
	conn, err := pgxpool.New(ctx, cfg.DSN())

	if err != nil {
		return nil, err
	}

	return &Client{
		db:     pg(conn),
		config: *cfg,
	}, nil
}

func (c *Client) DB() DB {
	return c.db
}

func (c *Client) Ping(ctx context.Context) error {
	return c.db.Ping(ctx)
}

func (c *Client) Close() error {
	return c.db.Close()
}
