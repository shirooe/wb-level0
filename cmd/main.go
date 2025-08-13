package main

import (
	"context"
	"log"
	"wb-level0/internal/kafka"
)

const topic = "order-topic"

func main() {
	ctx := context.Background()
	cfg := kafka.NewConfig([]string{"localhost:9092"}, topic, topic, 0, 1000000)

	c := kafka.NewConsumer(ctx, cfg)
	defer func() {
		if err := c.Close(); err != nil {
			log.Fatalf("failed to close consumer: %v", err)
		}
	}()

	for {
		msg, err := c.Read()

		if err != nil {
			log.Printf("failed to read message: %v", err)
			continue
		}

		if err := c.Commit(); err != nil {
			log.Printf("failed to commit message: %v", err)
		}

		log.Printf("message: %s", string(msg.Value))
	}
}
