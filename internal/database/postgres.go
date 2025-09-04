package database

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// ключ для транзакции
type key string

const TxKey key = "tx"

// проверка реалзиации всех функции
var _ DB = (*postgres)(nil)

// postgres структура БД
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

// выполнение запроса и сканирование в одну структуру
func (p *postgres) ScanOneContext(ctx context.Context, dest any, query Query, args ...any) error {
	rows, err := p.QueryContext(ctx, query, args...)
	if err != nil {
		p.log.Info("[psql] ошибка", zap.Error(err))
		return err
	}

	return pgxscan.ScanOne(dest, rows)
}

// выполнение запроса и сканирование в слайс структур
func (p *postgres) ScanAllContext(ctx context.Context, dest any, query Query, args ...any) error {
	rows, err := p.QueryContext(ctx, query, args...)
	if err != nil {
		p.log.Info("[psql] ошибка", zap.Error(err))
		return err
	}
	return pgxscan.ScanAll(dest, rows)
}

// выполнение запроса и возврат pgx.Row для дальнейшего взаимодействия
func (p *postgres) QueryRowContext(ctx context.Context, query Query, args ...any) pgx.Row {
	// проверка на существование транзакции
	tx, ok := ctx.Value(TxKey).(pgx.Tx)

	if ok {
		// выполнение запроса в рамках транзакции
		return tx.QueryRow(ctx, query.QueryRaw, args...)
	}

	// выполнение запроса
	return p.pool.QueryRow(ctx, query.QueryRaw, args...)
}

// выполнение запроса и возврат pgx.Rows для дальнейшего взаимодействия
func (p *postgres) QueryContext(ctx context.Context, query Query, args ...any) (pgx.Rows, error) {
	// проверка на существование транзакции
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		// выполнение запроса в рамках транзакции
		return tx.Query(ctx, query.QueryRaw, args...)
	}

	// выполнение запроса
	return p.pool.Query(ctx, query.QueryRaw, args...)
}

// выполнение запроса и отслеживание изменении через rows.Affected()
func (p *postgres) ExecContext(ctx context.Context, query Query, args ...any) (pgconn.CommandTag, error) {
	// проверка на существование транзакции
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		// выполнение запроса в рамках транзакции
		return tx.Exec(ctx, query.QueryRaw, args...)
	}

	// выполнение запроса
	return p.pool.Exec(ctx, query.QueryRaw, args...)
}

// начать выполнение транзакции
func (p *postgres) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return p.pool.BeginTx(ctx, txOptions)
}
