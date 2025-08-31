package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB interface {
	SQLExecer
	Transactor

	Pooler
	Pinger
	Closer
}

type Pooler interface {
	Pool() *pgxpool.Pool
}

type Pinger interface {
	Ping(ctx context.Context) error
}

type Closer interface {
	Close() error
}

type Query struct {
	Name     string
	QueryRaw string
}

type SQLExecer interface {
	NamedExecer
	QueryExecer
}

type NamedExecer interface {
	ScanOneContext(ctx context.Context, dest any, query Query, args ...any) error
	ScanAllContext(ctx context.Context, dest any, query Query, args ...any) error
}

type QueryExecer interface {
	QueryContext(ctx context.Context, query Query, args ...any) (pgx.Rows, error)
	QueryRowContext(ctx context.Context, query Query, args ...any) pgx.Row
	ExecContext(ctx context.Context, query Query, args ...any) (pgconn.CommandTag, error)
}

type Transactor interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

type Handler func(ctx context.Context) error
