package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"golang-kafka/internal/kafka"
	"golang-kafka/internal/order"
)

func main() {
	brokers := []string{"localhost:9092"}
	topic := "orders"
	groupID := "order-analytics-group"

	// Initialize consumer
	consumer := kafka.NewConsumer(brokers, topic, groupID)
	defer consumer.Close()

	// Create analyzer
	analyzer := order.NewAnalyzer(consumer)

	// Create context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start analyzer in goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- analyzer.Start(ctx)
	}()

	fmt.Println("Order analyzer started. Press Ctrl+C to exit.")

	select {
	case err := <-errChan:
		if err != nil {
			log.Printf("Error in analyzer: %v\n", err)
		}
	case <-sigChan:
		fmt.Println("\nReceived shutdown signal")
		cancel()
	}

	fmt.Println("Shutting down analyzer...")
}
