package main

import (
	"context"
	"golang-kafka/internal/kafka"
	"log"
)

func main() {
	brokers := []string{"localhost:9092"}
	topic := "message"
	groupID := "messages"

	consumer := kafka.NewConsumer(brokers, topic, groupID)
	defer consumer.Close()

	ctx := context.Background()

	for {
		key, value, err := consumer.Consume(ctx)
		if err != nil {
			log.Printf("Failed to consume message : %v\n", err)
			continue
		}

		log.Printf("Continue Message : key=%s   value=%s", key, value)
	}

}
