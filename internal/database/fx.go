package database

import (
	"context"
	"wb-level0/internal/database/migrations"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// регистрация модуля database
func Module() fx.Option {
	return fx.Module("database",
		// регистрация конфига и клиента БД
		fx.Provide(ProvideConfig, ProvideClient),
		// создание экземпляра структуры
		fx.Supply(&migrations.Migrations{}),
		// запуск модуля
		fx.Invoke(func(lc fx.Lifecycle, client *Client, log *zap.Logger, migration *migrations.Migrations) {
			// управления жизненным циклом модуля
			lc.Append(fx.Hook{
				// при старте
				OnStart: func(ctx context.Context) error {
					// пинг БД, проверка на работоспособность
					if err := client.Ping(ctx); err != nil {
						return err
					}

					// запуск миграции, в случае отсутствия таблиц в БД
					if err := migration.Up(client.config.DSN()); err != nil {
						return err
					}

					log.Info("[database] подключение к базе данных установлено")
					return nil
				},
				// при остановке
				OnStop: func(ctx context.Context) error {
					log.Info("[database] подключение к базе данных закрыто")
					// закрыть соединения к БД
					return client.Close()
				},
			})
		}),
	)
}
