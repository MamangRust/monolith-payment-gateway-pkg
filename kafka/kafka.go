package kafka

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
)

type Kafka struct {
	logger   logger.LoggerInterface
	producer sarama.SyncProducer
	brokers  []string
}

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

func (k *Kafka) GetBrokers() []string {
	return k.brokers
}

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
		for {
			if err := consumerGroup.Consume(ctx, topics, handler); err != nil {
				log.Printf("Error from consumer: %v", err)
			}
		}
	}()

	go func() {
		for err := range consumerGroup.Errors() {
			log.Printf("Consumer group error: %v", err)
		}
	}()

	return nil
}
