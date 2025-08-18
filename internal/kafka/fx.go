package kafka

import (
	"context"

	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("kafka",
		fx.Provide(ProvideConfig, ProvideConsumer),
		fx.Invoke(func(lc fx.Lifecycle, consumer *Consumer) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go consumer.Consume()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					return consumer.Close()
				},
			})
		}),
	)
}
