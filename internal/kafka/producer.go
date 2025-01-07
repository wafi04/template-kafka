package kafka

import (
	"context"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:          brokers,
		Topic:            topic,
		Async:            true,
		BatchSize:        100,
		RequiredAcks:     int(kafka.RequireOne),
		BatchTimeout:     time.Millisecond,
		CompressionCodec: kafka.Compression.Codec(kafka.Lz4),
	})

	return &Producer{
		writer: writer,
	}
}

func (p *Producer) Produce(ctx context.Context, key, value []byte) error {
	message := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}
	return p.writer.WriteMessages(ctx, message)
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
