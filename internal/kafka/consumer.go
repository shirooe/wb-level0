package kafka

import (
	"context"
	"fmt"
	"io"
	"log"
	"wb-level0/internal/service"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	ctx    context.Context
	reader *kafka.Reader

	service *service.WBLevel0Service
}

func ProvideConsumer(ctx context.Context, cfg *Config, service *service.WBLevel0Service) *Consumer {
	return &Consumer{
		ctx: ctx,
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

func (c *Consumer) Fetch() (kafka.Message, error) {
	return c.reader.FetchMessage(c.ctx)
}

func (c *Consumer) Consume() {
	for {
		msg, err := c.Fetch()
		if err != nil {
			if err == context.Canceled || err == io.EOF {
				// TODO: consumer was closed
				log.Printf("consumer stopped: %v", err)
				return
			}
			log.Printf("failed to read message: %v", err)
			continue
		}

		fmt.Printf("message: %s\n", msg.Value)

		if err := c.Commit(); err != nil {
			log.Printf("failed to commit message: %v", err)
		}
	}
}

func (c *Consumer) Commit() error {
	return c.reader.CommitMessages(c.ctx)
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
