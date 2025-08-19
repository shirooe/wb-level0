package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

func unmarshalToModel[T any](data []byte) (T, error) {
	var zero T
	if len(data) == 0 {
		return zero, errors.New("[util] пустой байтовой массив")
	}

	var model T
	err := json.Unmarshal(data, &model)
	if err != nil {
		return zero, err
	}
	return model, nil
}

func handlePgError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}

	switch pgErr.Code {
	case "23505":
		return fmt.Errorf("дублирующая запись: PG_ERROR_CODE [%s]", pgErr.Code)
	case "23503":
		return fmt.Errorf("несуществующая запись: PG_ERROR_CODE [%s]", pgErr.Code)
	default:
		return fmt.Errorf("ошибка базы данных: PG_ERROR_CODE [%s]", pgErr.Code)
	}
}
