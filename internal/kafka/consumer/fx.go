package consumer

import (
	"context"
	"wb-level0/internal/kafka/config"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// регистрация модуля
func Module() fx.Option {
	return fx.Module("kafka",
		// регистрация консьюмера
		fx.Provide(ProvideConsumer),
		// запуск консьюмера
		fx.Invoke(func(lc fx.Lifecycle, ctx context.Context, consumer *Consumer, cfg *config.Config, log *zap.Logger) {
			// управление жизненным циклом модуля
			lc.Append(fx.Hook{
				// при старте
				OnStart: func(_ context.Context) error {
					// создание топика для сообщении
					consumer.CreateTopic(cfg)
					// получение сообщении
					go consumer.Consume(ctx)
					log.Info("[kafka] консьюмер запущен")
					return nil
				},
				// при остановке
				OnStop: func(ctx context.Context) error {
					log.Info("[kafka] консьюмер остановлен")
					// остановление консьюмера
					return consumer.Close()
				},
			})
		}),
	)
}
