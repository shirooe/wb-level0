package database

import (
	"context"
	"wb-level0/internal/database/migrations"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Module() fx.Option {
	return fx.Module("database",
		fx.Provide(ProvideConfig, ProvideClient),
		fx.Supply(&migrations.Migrations{}),
		fx.Invoke(func(lc fx.Lifecycle, client *Client, log *zap.Logger, migration *migrations.Migrations) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					if err := client.Ping(ctx); err != nil {
						return err
					}

					if err := migration.Up(client.config.DSN()); err != nil {
						return err
					}

					log.Info("[database] подключение к базе данных установлено")
					return nil
				},
				OnStop: func(ctx context.Context) error {
					log.Info("[database] подключение к базе данных закрыто")
					return client.Close()
				},
			})
		}),
	)
}
