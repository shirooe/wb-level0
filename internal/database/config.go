package database

import (
	"fmt"
	"log"
	"os"

	"go.uber.org/config"
)

type Config struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func ProvideConfig(provider *config.YAML) *Config {
	var cfg Config

	if err := provider.Get("database").Populate(&cfg); err != nil {
		log.Fatalf("[Database] Ошибка конфигурации: %v", err)
	}

	return &cfg
}

func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
}
