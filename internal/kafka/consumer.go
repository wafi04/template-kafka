package kafka

import (
	"context"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MaxWait:  100 * time.Millisecond,
		MinBytes: 10e3, //  min 10kb
		MaxBytes: 10e6, // max 10mb
	})

	return &Consumer{
		reader: reader,
	}
}

func (c *Consumer) Consume(ctx context.Context) ([]byte, []byte, error) {
	// get message from context
	message, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return nil, nil, err
	}
	// mengembalikan key dan value
	return message.Key, message.Value, nil
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
