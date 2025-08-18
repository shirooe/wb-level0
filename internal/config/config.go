package config

import (
	"log"

	"github.com/joho/godotenv"
	"go.uber.org/config"
)

func New() *config.YAML {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("failed to load env: %v", err)
	}

	provider, err := config.NewYAML(config.File("config/config.yml"))

	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	return provider
}
