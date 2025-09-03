package consumer

import (
	"context"
	"wb-level0/internal/kafka/config"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Module() fx.Option {
	return fx.Module("kafka",
		fx.Provide(ProvideConsumer),
		fx.Invoke(func(lc fx.Lifecycle, ctx context.Context, consumer *Consumer, cfg *config.Config, log *zap.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					consumer.CreateTopic(cfg)
					go consumer.Consume(ctx)
					log.Info("[kafka] консьюмер запущен")
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
