package kafka

import (
	"context"
	"io"
	"net"
	"strconv"
	"wb-level0/internal/service"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Consumer struct {
	reader *kafka.Reader

	service *service.WBLevel0Service
	log     *zap.Logger
}

func ProvideConsumer(ctx context.Context, cfg *Config, service *service.WBLevel0Service) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:   cfg.Brokers,
			GroupID:   cfg.GroupID,
			Partition: cfg.Partition,
			Topic:     cfg.Topic,
			MaxBytes:  cfg.MaxBytes,
		}),
		service: service,
	}
}

func (c *Consumer) Fetch(ctx context.Context) (kafka.Message, error) {
	return c.reader.FetchMessage(ctx)
}

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

		c.service.CreateOrder(ctx, msg.Value)

		if err := c.Commit(ctx, msg); err != nil {
			c.log.Info("[kafka] ошибка коммита", zap.Error(err))
		}
	}
}

func (c *Consumer) Commit(ctx context.Context, msg kafka.Message) error {
	return c.reader.CommitMessages(ctx, msg)
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

func (c *Consumer) CreateTopic(config *Config) {
	conn, err := kafka.Dial("tcp", config.Brokers[0])
	if err != nil {
		c.log.Error("[kafka] ошибка подключения к брокеру", zap.Error(err))
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		c.log.Error("[kafka] ошибка получения контроллера", zap.Error(err))
	}

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		c.log.Error("[kafka] ошибка подключения к контроллеру", zap.Error(err))
	}
	defer controllerConn.Close()

	if err := controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             config.Topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}); err != nil {
		c.log.Error("[kafka] ошибка создания топика", zap.Error(err))
	}
}
