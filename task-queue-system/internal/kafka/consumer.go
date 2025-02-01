package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader *kafka.Reader
}

func NewKafkaConsumer(brokers []string, topic string) *KafkaConsumer {
	return &KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: "my-group",
		}),
	}
}

func (c *KafkaConsumer) ReadMessage(ctx context.Context) (string, error) {
	m, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return "", err
	}
	return string(m.Value), nil
}

func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}
