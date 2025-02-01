package taskqueue

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaTaskQueue struct {
	writer *kafka.Writer
	reader *kafka.Reader
}

func NewKafkaTaskQueue(brokers []string, topic string) *KafkaTaskQueue {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
	})

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  "my-group",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &KafkaTaskQueue{
		writer: writer,
		reader: reader,
	}
}

func (ktq *KafkaTaskQueue) Produce(ctx context.Context, topic string, message []byte) error {
	return ktq.writer.WriteMessages(ctx, kafka.Message{
		Value: message,
	})
}

func (ktq *KafkaTaskQueue) Consume(ctx context.Context, topic string, handler func([]byte) error) error {
	for {
		m, err := ktq.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}
		if err := handler(m.Value); err != nil {
			log.Printf("Failed to handle message: %v", err)
		}
	}
}

func (ktq *KafkaTaskQueue) Close() error {
	if err := ktq.writer.Close(); err != nil {
		return err
	}
	return ktq.reader.Close()
}
