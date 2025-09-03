package producer

import (
	"context"

	"wb-level0/internal/kafka/config"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Producer struct {
	writer *kafka.Writer
	log    *zap.Logger
}

func ProvideProducer(ctx context.Context, cfg *config.Config, log *zap.Logger) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:  kafka.TCP(cfg.Brokers...),
			Topic: cfg.Topic,
		},
		log: log,
	}
}

func (p *Producer) WriteTestMessage(data []byte) error {
	if err := p.writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte("test"),
		Value: data,
	}); err != nil {
		p.log.Info("[kafka] ошибка отправки сообщения", zap.Error(err))
		return err
	}

	p.log.Info("[kafka] тестовая модель отправлена")
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
