package http

import (
	"net"

	"go.uber.org/config"
)

type Config struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func ProvideConfig(provider *config.YAML) *Config {
	var cfg Config

	if err := provider.Get("server").Populate(&cfg); err != nil {
		panic(err)
	}

	return &cfg
}

func (c *Config) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}
