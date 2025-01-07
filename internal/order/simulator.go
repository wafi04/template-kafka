package order

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"golang-kafka/internal/kafka"
)

type Simulator struct {
	producer *kafka.Producer
}

func NewSimulator(producer *kafka.Producer) *Simulator {
	return &Simulator{
		producer: producer,
	}
}

func (s *Simulator) Start(ctx context.Context) error {
	products := []string{"Laptop", "Smartphone"}
	statuses := []string{"pending", "paid", "shipped", "delivered"}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			order := Order{
				OrderID:   generateID(),
				UserID:    generateID(),
				ProductID: products[rand.Intn(len(products))],
				Amount:    rand.Float64() * 1000,
				Status:    statuses[rand.Intn(len(statuses))],
				Timestamp: time.Now(),
			}

			orderJson, err := json.Marshal(order)
			if err != nil {
				return fmt.Errorf("error marshaling order: %w", err)
			}

			// Menggunakan producer yang sudah dikonfigurasi
			err = s.producer.Produce(ctx, []byte(order.OrderID), orderJson)
			if err != nil {
				return fmt.Errorf("error producing message: %w", err)
			}

			fmt.Printf("Produced order: %s\n", string(orderJson))
			time.Sleep(time.Second)
		}
	}
}

func generateID() string {
	return fmt.Sprintf("%d", rand.Int63())
}
