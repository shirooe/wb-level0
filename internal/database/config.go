package database

import (
	"fmt"
	"os"

	"go.uber.org/config"
	"go.uber.org/zap"
)

type Config struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// получение провайдера и маппинг данных в структуру из yml-файла по ключу database
func ProvideConfig(provider *config.YAML, log *zap.Logger) *Config {
	var cfg Config

	if err := provider.Get("database").Populate(&cfg); err != nil {
		log.Info("[database] ошибка конфигурации", zap.Error(err))
	}

	return &cfg
}

// DSN для подключения к БД
func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
}
