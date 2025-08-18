package kafka

import "go.uber.org/config"

type Config struct {
	Brokers   []string `yaml:"brokers"`
	Topic     string   `yaml:"topic"`
	GroupID   string   `yaml:"group_id"`
	Partition int      `yaml:"partition"`
	MaxBytes  int      `yaml:"max_bytes"`
}

func ProvideConfig(provider *config.YAML) *Config {
	var cfg Config

	if err := provider.Get("kafka").Populate(&cfg); err != nil {
		panic(err)
	}

	return &cfg
}
