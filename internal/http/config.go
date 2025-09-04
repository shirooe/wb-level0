package http

import (
	"net"

	"go.uber.org/config"
	"go.uber.org/zap"
)

type Config struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// получение провайдера и маппинг данных в структуру из yml-файла по ключу server
func ProvideConfig(provider *config.YAML, log *zap.Logger) *Config {
	var cfg Config

	if err := provider.Get("server").Populate(&cfg); err != nil {
		log.Info("[server] ошибка конфигурации", zap.Error(err))
	}

	return &cfg
}

// объединение Config.Host и Config.Port
func (c *Config) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}
