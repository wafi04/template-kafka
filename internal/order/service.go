package order

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/segmentio/kafka-go"
)

func generateID() string {
	return fmt.Sprintf("%d", rand.Int63())
}

func simulateOrders(ctx context.Context, writer *kafka.Writer) {
	products := []string{"Laptop", "Smartphone"}
	statuses := []string{"pending", "paid", "shipped", "delivered"}

	for {
		order := Order{
			OrderID:   generateID(),
			UserID:    generateID(),
			ProductID: products[rand.Intn(len(products))],
			Amount:    rand.Float64() * 1000,
			Status:    statuses[rand.Intn(len(statuses))],
			Timestamp: time.Now(),
		}

		orderJson, _ := json.Marshal(order)
		fmt.Println(string(orderJson))
		err := writer.WriteMessages(ctx, kafka.Message{
			Value: orderJson,
		})

		if err != nil {
			fmt.Printf("Error  Proccedd data %v\n", err)
		}
		time.Sleep(time.Second)

	}
}

func processAnalytics(ctx context.Context, reader *kafka.Reader) {
	analytics := Analytics{
		ProductCounts: make(map[string]int),
		HourlyRevenue: make(map[int]float64),
	}

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		var order Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Printf("Error unmarshaling order: %v", err)
			continue
		}

		// Update analytics
		analytics.TotalOrders++
		analytics.TotalRevenue += order.Amount
		analytics.ProductCounts[order.ProductID]++
		hour := order.Timestamp.Hour()
		analytics.HourlyRevenue[hour] += order.Amount

		// Log current analytics
		logAnalytics(analytics)
	}
}

func logAnalytics(a Analytics) {
	log.Printf("\n=== Current Analytics ===\n")
	log.Printf("Total Orders: %d\n", a.TotalOrders)
	log.Printf("Total Revenue: $%.2f\n", a.TotalRevenue)
	log.Printf("Product Counts: %v\n", a.ProductCounts)
	log.Printf("Hourly Revenue: %v\n", a.HourlyRevenue)
}
