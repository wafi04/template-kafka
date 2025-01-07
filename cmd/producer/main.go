package main

import (
	"bufio"
	"context"
	"fmt"
	"golang-kafka/internal/kafka"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	brokers := []string{"localhost:9092"}
	topic := "test_topic"

	producer := kafka.NewProducer(brokers, topic)
	defer producer.Close()

	ctx := context.Background()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter message (or 'quit' to exit): ")
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading input: %v\n", err)
			continue
		}

		message = strings.TrimSpace(message)

		if message == "quit" {
			break
		}

		key := []byte(fmt.Sprintf("key-%d", time.Now().Unix()))
		value := []byte(message)

		err = producer.Produce(ctx, key, value)
		if err != nil {
			log.Printf("Failed to produce message: %v\n", err)
			continue
		}

		log.Printf("Produced Message: key=%s value=%s\n", key, value)
	}
}
