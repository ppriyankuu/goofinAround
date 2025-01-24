package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

var Writer *kafka.Writer
var Reader *kafka.Reader

func InitProducer(brokers []string, topic string) {
	Writer = &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	log.Println("Kafka producer initialized.")
}

func ProduceMessage(key, value []byte) error {
	err := Writer.WriteMessages(context.Background(), kafka.Message{
		Key:   key,
		Value: value,
	})
	if err != nil {
		log.Printf("Failed to produce message: %v", err)
		return err
	}
	log.Println("Message produced to Kafka.")
	return nil
}

func InitConsumer(brokers []string, topic, groupID string) {
	Reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})
	log.Println("Kafka consumer initialized.")
}

func handler(key, value []byte) {
	log.Printf("Processing message with key: %s and value: %s", string(key), string(value))
}

func Consume(topic string) {
	for {
		msg, err := Reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error consuming message: %v", err)
			continue
		}
		log.Printf("Message consumed: key=%s value=%s", string(msg.Key), string(msg.Value))
		handler(msg.Key, msg.Value)
	}
}

func CloseKafka() {
	if Writer != nil {
		_ = Writer.Close()
	}
	if Reader != nil {
		_ = Reader.Close()
	}
	log.Println("Kafka connections closed.")
}
