package main

import (
	"context"
	"log"
	"wb-level0/internal/app"
)

// точка входа в приложение
func main() {
	// создание регистратора fx
	application := app.New()
	ctx := context.Background()

	// старт приложения
	if err := application.Start(ctx); err != nil {
		log.Fatalf("❌ Ошибка старта приложения: %v", err)
	}

	// канал для блокировки выполнения
	<-application.Done()

	// остановка приложения
	if err := application.Stop(ctx); err != nil {
		log.Fatalf("❌ Ошибка остановки приложения: %v", err)
	}
}
