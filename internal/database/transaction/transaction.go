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
	// проверка выполняется ли транзакция к текущий момент
	tx, ok := ctx.Value(database.TxKey).(pgx.Tx)
	// в случае выполнения пропустить повторное старт
	if ok {
		return fn(ctx)
	}

	// начало транзакции
	tx, err := m.db.DB().BeginTx(ctx, opts)
	if err != nil {
		m.log.Error("[psql] ошибка при создании транзакции", zap.Error(err))
		return errors.Join(err, fmt.Errorf("[psql] ошибка при создании транзакции %v", err))
	}

	// сохранение транзакции в контексте
	ctx = context.WithValue(ctx, database.TxKey, tx)

	defer func() {
		if err != nil {
			// отмена выполнения транзакции в случае ошибки внутри транзакции
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				err = errors.Join(err, fmt.Errorf("[psql] ошибка при откате транзакции %v", errRollback))
				m.log.Error("[psql] ошибка при откате транзакции", zap.Error(errRollback))
			}
		}
	}()

	// выполнение функции
	if err := fn(ctx); err != nil {
		return err
	}

	// зафиксировать изменение в БД после успешнего выполнения транзакции
	if err := tx.Commit(ctx); err != nil {
		m.log.Error("[psql] ошибка при завершении транзакции", zap.Error(err))
		err = errors.Join(err, fmt.Errorf("[psql] ошибка при завершении транзакции %v", err))
		return err
	}

	return nil
}

// выполнение транзакции с уровнем изоляции Read Committed
func (m *Manager) WithTransaction(ctx context.Context, fn database.Handler) error {
	return m.transaction(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, fn)
}
