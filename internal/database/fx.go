package database

import (
	"context"

	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("database",
		fx.Provide(ProvideConfig, ProvideClient),
		fx.Invoke(func(lc fx.Lifecycle, client *Client) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					if err := client.Ping(ctx); err != nil {
						return err
					}
					return nil
				},
				OnStop: func(ctx context.Context) error {
					return client.Close()
				},
			})
		}),
	)
}
