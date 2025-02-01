package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	return &KafkaProducer{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers: brokers,
			Topic:   topic,
		}),
	}
}

func (p *KafkaProducer) SendMessage(ctx context.Context, message string) error {
	err := p.writer.WriteMessages(ctx, kafka.Message{
		Value: []byte(message),
	})
	if err != nil {
		// Logs the error or wraps it with additional context
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
