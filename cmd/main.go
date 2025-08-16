package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"wb-level0/internal/database"
	"wb-level0/internal/kafka"
	"wb-level0/internal/models"
	"wb-level0/internal/repository"

	"github.com/joho/godotenv"
)

const topic = "order-topic"

func main() {
	ctx := context.Background()
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load env: %v", err)
	}

	cfg := kafka.NewConfig([]string{"kafka:9092"}, topic, topic, 0, 1000000)

	dbCfg := database.NewConfig("localhost", "5432")
	db, err := database.New(ctx, dbCfg)
	if err != nil {
		log.Fatalf("failed to create database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("failed to close database: %v", err)
		}
	}()

	orderRepository := repository.NewOrderRepository(db)
	deliveryRepository := repository.NewDeliveryRepository(db)
	itemRepository := repository.NewItemRepository(db)
	paymentRepository := repository.NewPaymentRepository(db)

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

		var v models.Order
		json.Unmarshal(msg.Value, &v)

		fmt.Println("message received")

		if _, err := orderRepository.Create(ctx, v); err != nil {
			log.Printf("failed to create order: %v", err)
		}

		if _, err := deliveryRepository.Create(ctx, v.OrderUID, v.Delivery); err != nil {
			log.Printf("failed to create delivery: %v", err)
		}

		if _, err := paymentRepository.Create(ctx, v.OrderUID, v.Payment); err != nil {
			log.Printf("failed to create payment: %v", err)
		}

		if _, err := itemRepository.Create(ctx, v.OrderUID, v.Items); err != nil {
			log.Printf("failed to create items: %v", err)
		}

		fmt.Println("message processed")

		if err := c.Commit(); err != nil {
			log.Printf("failed to commit message: %v", err)
		}
	}
}
