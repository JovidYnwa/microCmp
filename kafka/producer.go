package kafka

// import (
// 	"log"

// )

// type KafkaProducer struct {
// 	producer *kafka.
// }

// func NewKafkaProducer(kafkaServerIP string) (*KafkaProducer, error) {
// 	config := &kafka.ConfigMap{
// 		"bootstrap.servers": kafkaServerIP, // IP address of the Kafka server
// 	}

// 	producer, err := kafka.NewProducer(config)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &KafkaProducer{producer: producer}, nil
// }

// func (kp *KafkaProducer) SendMessage(topic string, message string) error {
// 	// Create the Kafka message
// 	msg := &kafka.Message{
// 		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
// 		Value:          []byte(message),
// 	}

// 	// Produce the message to the topic
// 	err := kp.producer.Produce(msg, nil)
// 	if err != nil {
// 		log.Printf("Error sending message to Kafka: %v", err)
// 		return err
// 	}
// 	log.Printf("Message sent to Kafka topic %s: %s", topic, message)
// 	return nil
// }

// func (kp *KafkaProducer) Close() {
// 	kp.producer.Close()
// }
