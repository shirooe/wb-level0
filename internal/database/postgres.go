package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgres struct {
	pool *pgxpool.Pool
}

func pg(conn *pgxpool.Pool) *postgres {
	return &postgres{
		pool: conn,
	}
}

func (p postgres) Pool() *pgxpool.Pool {
	return p.pool
}

func (p postgres) Ping(ctx context.Context) error {
	return p.pool.Ping(ctx)
}

func (p postgres) Close() error {
	p.pool.Close()
	return nil
}

func (p *postgres) ScanOneContext(ctx context.Context, dest interface{}, query Query, args ...interface{}) error {
	// rows, err := p.QueryContext(ctx, query, args...)
	// if err != nil {
	// log.Error("failed to query %v", err)
	// return err
	// }
	// return pgxscan.ScanOne(dest, rows)
	return nil
}

func (p *postgres) ScanAllContext(ctx context.Context, dest interface{}, query Query, args ...interface{}) error {
	// rows, err := p.QueryContext(ctx, query, args...)
	// if err != nil {
	// 	log.Error("failed to query %v", err)
	// 	return err
	// }
	// return pgxscan.ScanAll(dest, rows)
	return nil
}

func (p *postgres) QueryRowContext(ctx context.Context, query Query, args ...interface{}) pgx.Row {
	return p.pool.QueryRow(ctx, query.QueryRaw, args...)
}
