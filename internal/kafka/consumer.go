package kafka

import (
	"context"
	"io"
	"log"
	"wb-level0/internal/service"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader

	service *service.WBLevel0Service
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
				log.Println("[kafka] канал закрыт")
				return
			}
			log.Printf("[kafka] ошибка получения сообщения %v\n", err)
			continue
		}

		c.service.CreateOrder(ctx, msg.Value)

		if err := c.Commit(ctx); err != nil {
			log.Printf("[kafka] ошибка коммита %v\n", err)
		}
	}
}

func (c *Consumer) Commit(ctx context.Context) error {
	return c.reader.CommitMessages(ctx)
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
