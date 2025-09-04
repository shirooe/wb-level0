package config

import (
	"github.com/joho/godotenv"
	"go.uber.org/config"
	"go.uber.org/zap"
)

// создание конфига
func New(log *zap.Logger) *config.YAML {
	// загрузка .env-файла
	if err := godotenv.Load(".env"); err != nil {
		log.Info("[config] ошибка загрузки файла", zap.Error(err))
	}

	// загрузка провайдера для yml-файла
	provider, err := config.NewYAML(config.File("config/config.yml"))

	if err != nil {
		log.Info("[config] ошибка конфигурации", zap.Error(err))
	}

	// возвращение провайдера для считывания значении по ключу в yml-файле
	return provider
}
