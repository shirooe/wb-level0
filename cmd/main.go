package main

import (
	"context"
	"log"
	"wb-level0/internal/app"
)

func main() {
	application := app.New()
	ctx := context.Background()

	if err := application.Start(ctx); err != nil {
		log.Fatalf("❌ Ошибка старта приложения: %v", err)
	}

	<-application.Done()

	if err := application.Stop(ctx); err != nil {
		log.Fatalf("❌ Ошибка остановки приложения: %v", err)
	}
}
