package consumer

import (
	"context"
	"io"
	"net"
	"strconv"
	"wb-level0/internal/kafka/config"
	"wb-level0/internal/service"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Consumer struct {
	reader *kafka.Reader

	service *service.WBLevel0Service
	log     *zap.Logger
}

// создание консьюмера
func ProvideConsumer(ctx context.Context, cfg *config.Config, service *service.WBLevel0Service, log *zap.Logger) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:   cfg.Brokers,
			GroupID:   cfg.GroupID,
			Partition: cfg.Partition,
			Topic:     cfg.Topic,
			MaxBytes:  cfg.MaxBytes,
		}),
		service: service,
		log:     log,
	}
}

// получение сообщении
func (c *Consumer) Fetch(ctx context.Context) (kafka.Message, error) {
	return c.reader.FetchMessage(ctx)
}

// получение и коммит сообщении
func (c *Consumer) Consume(ctx context.Context) {
	for {
		msg, err := c.Fetch(ctx)
		if err != nil {
			if err == context.Canceled || err == io.EOF {
				// TODO: consumer was closed
				return
			}
			c.log.Info("[kafka] ошибка получения сообщения", zap.Error(err))
			continue
		}

		// создание заказа в сервисе
		c.service.CreateOrder(ctx, msg.Value)

		if err := c.Commit(ctx, msg); err != nil {
			c.log.Info("[kafka] ошибка коммита", zap.Error(err))
		}
	}
}

// коммит сообщении
func (c *Consumer) Commit(ctx context.Context, msg kafka.Message) error {
	return c.reader.CommitMessages(ctx, msg)
}

// закрытие соединении
func (c *Consumer) Close() error {
	return c.reader.Close()
}

// создание топика
func (c *Consumer) CreateTopic(config *config.Config) {
	// соединение с кафкой
	conn, err := kafka.Dial("tcp", config.Brokers[0])
	if err != nil {
		c.log.Error("[kafka] ошибка подключения к брокеру", zap.Error(err))
	}
	defer conn.Close()

	// получение контроллера
	controller, err := conn.Controller()
	if err != nil {
		c.log.Error("[kafka] ошибка получения контроллера", zap.Error(err))
	}

	// получение соединении контроллера
	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		c.log.Error("[kafka] ошибка подключения к контроллеру", zap.Error(err))
	}
	defer controllerConn.Close()

	// создание топика
	if err := controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             config.Topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}); err != nil {
		c.log.Error("[kafka] ошибка создания топика", zap.Error(err))
	}
}
