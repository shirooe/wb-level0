package transaction

import (
	"context"
	"errors"
	"fmt"
	"wb-level0/internal/database"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func (m *Manager) transaction(ctx context.Context, opts pgx.TxOptions, fn database.Handler) error {
	tx, ok := ctx.Value(database.TxKey).(pgx.Tx)
	if ok {
		return fn(ctx)
	}

	tx, err := m.db.DB().BeginTx(ctx, opts)
	if err != nil {
		m.log.Error("[psql] ошибка при создании транзакции", zap.Error(err))
		return errors.Join(err, fmt.Errorf("[psql] ошибка при создании транзакции %v", err))
	}

	ctx = context.WithValue(ctx, database.TxKey, tx)

	defer func() {
		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				err = errors.Join(err, fmt.Errorf("[psql] ошибка при откате транзакции %v", errRollback))
				m.log.Error("[psql] ошибка при откате транзакции", zap.Error(errRollback))
			}
		}
	}()

	if err := fn(ctx); err != nil {
		m.log.Error("[psql] ошибка при выполнении транзакции", zap.Error(err))
		err = errors.Join(err, fmt.Errorf("[psql] ошибка при выполнении транзакции %v", err))
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		m.log.Error("[psql] ошибка при завершении транзакции", zap.Error(err))
		err = errors.Join(err, fmt.Errorf("[psql] ошибка при завершении транзакции %v", err))
		return err
	}

	return nil
}

func (m *Manager) WithTransaction(ctx context.Context, fn database.Handler) error {
	return m.transaction(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, fn)
}
