package kafka

import (
	"context"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
)

// SyncProducer is an interface that represents a Kafka producer.
type SyncProducer interface {
	SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error)
	Close() error
}

// Kafka is a struct that represents a Kafka producer.
type Kafka struct {
	logger   logger.LoggerInterface
	producer SyncProducer
	brokers  []string
}

// NewKafka initializes a new Kafka struct.
//
// It takes a logger and a list of broker addresses as inputs and returns a pointer to the Kafka struct.
// It creates a new Kafka producer with the given configuration, and logs a message indicating if the connection is successful.
// If the connection fails, it logs an error message and exits.
func NewKafka(logger logger.LoggerInterface, brokers []string) *Kafka {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}

	log.Println("Kafka producer connected successfully")

	return &Kafka{
		producer: producer,
		brokers:  brokers,
		logger:   logger,
	}
}

// GetBrokers returns a list of the Kafka broker addresses that the producer is connected to.
func (k *Kafka) GetBrokers() []string {
	return k.brokers
}

// SendMessage sends a message to the given Kafka topic with the given key and value.
//
// It uses the configured SyncProducer to send the message and logs the result of the send operation.
// If the send operation fails, it returns an error.
func (k *Kafka) SendMessage(topic string, key string, value []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	partition, offset, err := k.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	return nil
}

// StartConsumers starts a Kafka consumer group to consume messages from the specified topics.
//
// It takes a list of topics, a consumer group ID, and a handler implementing the sarama.ConsumerGroupHandler interface.
// The method initializes a new Kafka consumer group and begins consuming messages in a background goroutine.
// If an error occurs during consumption, it retries up to a maximum number of retries with a delay between attempts.
// Any errors from the consumer group are logged and the function returns an error if the consumer group initialization fails.
func (k *Kafka) StartConsumers(topics []string, groupID string, handler sarama.ConsumerGroupHandler) error {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumerGroup, err := sarama.NewConsumerGroup(k.brokers, groupID, config)
	if err != nil {
		return err
	}

	ctx := context.Background()

	go func() {
		retries := 0
		maxRetries := 5
		for {
			err := consumerGroup.Consume(ctx, topics, handler)
			if err != nil {
				log.Printf("Error from consumer: %v", err)
				retries++
				if retries >= maxRetries {
					log.Fatalf("Max retries reached for consumer group. Exiting.")
				}
				time.Sleep(30 * time.Second)
				continue
			}
			retries = 0
		}
	}()

	go func() {
		for err := range consumerGroup.Errors() {
			log.Printf("Consumer group error: %v", err)
		}
	}()

	return nil
}
