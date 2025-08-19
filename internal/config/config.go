package config

import (
	"log"

	"github.com/joho/godotenv"
	"go.uber.org/config"
)

func New() *config.YAML {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("[config] ошибка загрузки файла %v", err)
	}

	provider, err := config.NewYAML(config.File("config/config.yml"))

	if err != nil {
		log.Fatalf("[config] ошибка конфигурации %v", err)
	}

	return provider
}
