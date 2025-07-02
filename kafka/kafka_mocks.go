package kafka

import (
	"github.com/IBM/sarama"
)

type MockProducer struct {
	Messages   []*sarama.ProducerMessage
	ShouldFail bool
}

func (m *MockProducer) SendMessage(msg *sarama.ProducerMessage) (int32, int64, error) {
	if m.ShouldFail {
		return 0, 0, sarama.ErrUnknown
	}
	m.Messages = append(m.Messages, msg)
	return 0, 123, nil
}

func (m *MockProducer) Close() error {
	return nil
}
