package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

type Message struct {
	Key   string
	Value interface{}
}

func NewProducer(brokers []string, topic string) *Producer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
		Async:   true,
		Logger: kafka.LoggerFunc(func(msg string, args ...interface{}) {
			log.Printf(msg, args...)
		}),
	})

	return &Producer{
		writer: writer,
	}
}

func (p *Producer) SendMessage(ctx context.Context, msg Message) error {
	value, err := json.Marshal(msg.Value)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(msg.Key),
		Value: value,
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

type ProducerConfig struct {
	Brokers []string
	Topic   string
}

func NewProducerFromEnv() *Producer {
	config := ProducerConfig{
		Brokers: []string{os.Getenv("KAFKA_BROKER")},
		Topic:   os.Getenv("KAFKA_TOPIC"),
	}

	return NewProducer(config.Brokers, config.Topic)
}
