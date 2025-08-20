package kafka

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Module() fx.Option {
	return fx.Module("kafka",
		fx.Provide(ProvideConfig, ProvideConsumer),
		fx.Invoke(func(lc fx.Lifecycle, ctx context.Context, consumer *Consumer, log *zap.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					go consumer.Consume(ctx)
					return nil
				},
				OnStop: func(ctx context.Context) error {
					log.Info("[kafka] консьюмер остановлен")
					return consumer.Close()
				},
			})
		}),
	)
}
