package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	ctx    context.Context
	reader *kafka.Reader
}

func NewConsumer(ctx context.Context, cfg *Config) *Consumer {
	return &Consumer{
		ctx: ctx,
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:   cfg.Brokers,
			GroupID:   cfg.GroupID,
			Partition: cfg.Partition,
			Topic:     cfg.Topic,
			MaxBytes:  cfg.MaxBytes,
		}),
	}
}

func (c *Consumer) Read() (kafka.Message, error) {
	return c.reader.ReadMessage(c.ctx)
}

func (c *Consumer) Commit() error {
	return c.reader.CommitMessages(c.ctx)
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
