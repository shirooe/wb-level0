package kafka

import (
	"context"

	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("kafka",
		fx.Provide(ProvideConfig, ProvideConsumer),
		fx.Invoke(func(lc fx.Lifecycle, ctx context.Context, consumer *Consumer) {
			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					go consumer.Consume(ctx)
					return nil
				},
				OnStop: func(ctx context.Context) error {
					return consumer.Close()
				},
			})
		}),
	)
}
