package kafka

import (
	"testing"

	mock_logger "github.com/MamangRust/monolith-payment-gateway-pkg/logger/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestKafka_SendMessage_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockProducer := &MockProducer{}

	k := &Kafka{
		producer: mockProducer,
		brokers:  []string{"localhost:9092"},
		logger:   mockLogger,
	}

	err := k.SendMessage("test-topic", "key", []byte("hello"))
	assert.NoError(t, err)
	assert.Len(t, mockProducer.Messages, 1)
	assert.Equal(t, "test-topic", mockProducer.Messages[0].Topic)
}

func TestKafka_SendMessage_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mock_logger.NewMockLoggerInterface(ctrl)
	mockProducer := &MockProducer{ShouldFail: true}

	k := &Kafka{
		producer: mockProducer,
		brokers:  []string{"localhost:9092"},
		logger:   mockLogger,
	}

	err := k.SendMessage("test-topic", "key", []byte("fail"))
	assert.Error(t, err)
}

func TestKafka_GetBrokers(t *testing.T) {
	k := &Kafka{
		brokers: []string{"broker1", "broker2"},
	}

	assert.Equal(t, []string{"broker1", "broker2"}, k.GetBrokers())
}
