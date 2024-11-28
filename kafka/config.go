package kafka

import "github.com/confluentinc/confluent-kafka-go/kafka"

type Config struct {
	BootstrapServers string
	Topic            string
	GroupID          string
	d                *kafka.Producer
}
