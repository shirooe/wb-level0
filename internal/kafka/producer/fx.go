package producer

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// регистрация модуля
func Module() fx.Option {
	return fx.Module(
		"producer",
		// регистрация продьюсера
		fx.Provide(ProvideProducer),
		// запуск модуля
		fx.Invoke(func(lc fx.Lifecycle, producer *Producer, log *zap.Logger) {
			// управление жизненным циклом модуля
			lc.Append(fx.Hook{
				// при старте
				OnStart: func(ctx context.Context) error {
					log.Info("[kafka] продюсер запущен")
					return nil
				},
				// при остановке
				OnStop: func(ctx context.Context) error {
					log.Info("[kafka] продюсер остановлен")
					// закрыть соединение
					return producer.writer.Close()
				},
			})
		}),
	)
}
