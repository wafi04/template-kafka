package main

import (
	"context"
	"fmt"
	"golang-kafka/internal/kafka"
	"golang-kafka/internal/order"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	brokers := []string{"localhost:9092"}
	topic := "orders"

	// Menggunakan producer yang sudah dikonfigurasi
	producer := kafka.NewProducer(brokers, topic)
	defer producer.Close()

	// Membuat simulator dengan producer yang sudah dikonfigurasi
	simulator := order.NewSimulator(producer)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start  simulator in goroutines
	errChan := make(chan error, 1)
	go func() {
		errChan <- simulator.Start(ctx)
	}()

	fmt.Println("Order simulator started. Press Ctrl+C to exit.")

	// Wait for error or shutdown signal
	select {
	case err := <-errChan:
		if err != nil {
			log.Printf("Error in simulator: %v\n", err)
		}
	case <-sigChan:
		fmt.Println("\nReceived shutdown signal")
		cancel()
	}

	fmt.Println("Shutting down simulator...")
}
