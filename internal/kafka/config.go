package kafka

type Config struct {
	Brokers   []string
	Topic     string
	GroupID   string
	Partition int
	MaxBytes  int
}

func NewConfig(brokers []string, topic string, groupID string, partition int, maxBytes int) *Config {
	return &Config{
		Brokers:   brokers,
		Topic:     topic,
		GroupID:   groupID,
		Partition: partition,
		MaxBytes:  maxBytes,
	}
}
