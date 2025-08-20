package config

import (
	"github.com/joho/godotenv"
	"go.uber.org/config"
	"go.uber.org/zap"
)

func New(log *zap.Logger) *config.YAML {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("[config] ошибка загрузки файла", zap.Error(err))
	}

	provider, err := config.NewYAML(config.File("config/config.yml"))

	if err != nil {
		log.Fatal("[config] ошибка конфигурации", zap.Error(err))
	}

	return provider
}
