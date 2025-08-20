package kafka

import (
	"go.uber.org/config"
	"go.uber.org/zap"
)

type Config struct {
	Brokers   []string `yaml:"brokers"`
	Topic     string   `yaml:"topic"`
	GroupID   string   `yaml:"group_id"`
	Partition int      `yaml:"partition"`
	MaxBytes  int      `yaml:"max_bytes"`
}

func ProvideConfig(provider *config.YAML, log *zap.Logger) *Config {
	var cfg Config

	if err := provider.Get("kafka").Populate(&cfg); err != nil {
		log.Fatal("[kafka] ошибка конфигурации", zap.Error(err))
	}

	return &cfg
}
