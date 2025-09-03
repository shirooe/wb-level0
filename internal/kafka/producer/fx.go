package producer

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Module() fx.Option {
	return fx.Module(
		"producer",
		fx.Provide(ProvideProducer),
		fx.Invoke(func(lc fx.Lifecycle, producer *Producer, log *zap.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					log.Info("[kafka] продюсер запущен")
					return nil
				},
				OnStop: func(ctx context.Context) error {
					log.Info("[kafka] продюсер остановлен")
					return producer.writer.Close()
				},
			})
		}),
	)
}
