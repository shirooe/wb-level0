package database

import (
	"fmt"
	"os"
)

type Config struct {
	Host string
	Port string
}

func NewConfig(host, port string) *Config {
	return &Config{
		Host: host,
		Port: port,
	}
}

func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
}
