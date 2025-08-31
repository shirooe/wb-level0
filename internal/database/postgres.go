package database

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type key string

const TxKey key = "tx"

var _ DB = (*postgres)(nil)

type postgres struct {
	pool *pgxpool.Pool
	log  *zap.Logger
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

func (p *postgres) ScanOneContext(ctx context.Context, dest any, query Query, args ...any) error {
	rows, err := p.QueryContext(ctx, query, args...)
	if err != nil {
		p.log.Info("[psql] ошибка", zap.Error(err))
		return err
	}

	return pgxscan.ScanOne(dest, rows)
}

func (p *postgres) ScanAllContext(ctx context.Context, dest any, query Query, args ...any) error {
	rows, err := p.QueryContext(ctx, query, args...)
	if err != nil {
		p.log.Info("[psql] ошибка", zap.Error(err))
		return err
	}
	return pgxscan.ScanAll(dest, rows)
}

func (p *postgres) QueryRowContext(ctx context.Context, query Query, args ...any) pgx.Row {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)

	if ok {
		return tx.QueryRow(ctx, query.QueryRaw, args...)
	}

	return p.pool.QueryRow(ctx, query.QueryRaw, args...)
}

func (p *postgres) QueryContext(ctx context.Context, query Query, args ...any) (pgx.Rows, error) {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Query(ctx, query.QueryRaw, args...)
	}

	return p.pool.Query(ctx, query.QueryRaw, args...)
}

func (p *postgres) ExecContext(ctx context.Context, query Query, args ...any) (pgconn.CommandTag, error) {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Exec(ctx, query.QueryRaw, args...)
	}

	return p.pool.Exec(ctx, query.QueryRaw, args...)
}

func (p *postgres) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return p.pool.BeginTx(ctx, txOptions)
}
