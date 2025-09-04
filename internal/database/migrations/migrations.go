package migrations

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migrations struct{}

// начать выполнение миграции
func (m *Migrations) Up(dsn string) error {
	// открыть соединение
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	// закрыть соединение
	defer db.Close()

	// получение драйвера postgres
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	// закрытие драйвера
	defer driver.Close()

	// получение инстанса миграции для взаимодействия
	migrator, err := migrate.NewWithDatabaseInstance("file://migrations/", "postgres", driver)
	if err != nil {
		return err
	}
	// закрытие инстанса
	defer migrator.Close()

	// выполнение Up функции
	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

// откат миграции
func (m *Migrations) Down(dsn string) error {
	// открыть соединение
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	// закрыть соединение
	defer db.Close()

	// получение драйвера postgres
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	// закрытие драйвера
	defer driver.Close()

	// получение инстанса миграции для взаимодействия
	migrator, err := migrate.NewWithDatabaseInstance("file://migrations/", "postgres", driver)
	if err != nil {
		return err
	}
	// закрытие инстанса
	defer migrator.Close()

	// выполнение Down функции
	if err := migrator.Down(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
