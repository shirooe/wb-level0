package kafka

import (
	"wb-level0/internal/kafka/config"
	"wb-level0/internal/kafka/consumer"
	"wb-level0/internal/kafka/producer"

	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		consumer.Module(),
		producer.Module(),
		fx.Provide(config.ProvideConfig),
	)
}
