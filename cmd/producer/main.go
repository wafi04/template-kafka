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
	topic := "order"

	// Inisialisasi producer
	producer := kafka.NewProducer(brokers, topic)
	defer producer.Close()

	ctx := context.Background()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Kafka Producer dimulai. Ketik 'quit' untuk keluar.")

	errChan := make(chan error, 1)
	go handleErrors(errChan)

	for {
		// Baca input
		fmt.Print("Masukkan pesan: ")
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error membaca input: %v\n", err)
			continue
		}

		message = strings.TrimSpace(message)

		// Periksa kondisi keluar
		if message == "quit" || message == "" {
			break
		}

		// Generate key dan value
		key := generateMessageKey()
		value := []byte(message)

		// Kirim pesan secara asynchronous
		go func(k, v []byte) {
			if err := producer.Produce(ctx, k, v); err != nil {
				errChan <- fmt.Errorf("gagal mengirim pesan: %v", err)
				return
			}
			log.Printf("Pesan terkirim: key=%s value=%s\n", k, v)
		}(key, value)
	}
}

// generateMessageKey menghasilkan key unik untuk setiap pesan
func generateMessageKey() []byte {
	return []byte(fmt.Sprintf("key-%d-%d", time.Now().Unix(), time.Now().Nanosecond()))
}

// handleErrors menangani error secara asynchronous
func handleErrors(errChan <-chan error) {
	for err := range errChan {
		log.Printf("Error: %v\n", err)
	}
}
