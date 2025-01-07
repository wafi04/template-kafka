package order

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"golang-kafka/internal/kafka"
)

type Analyzer struct {
	consumer  *kafka.Consumer
	analytics Analytics
}

func NewAnalyzer(consumer *kafka.Consumer) *Analyzer {
	return &Analyzer{
		consumer: consumer,
		analytics: Analytics{
			ProductCounts: make(map[string]int),
			HourlyRevenue: make(map[int]float64),
		},
	}
}

func (a *Analyzer) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// get key,value,err from consumer
			key, value, err := a.consumer.Consume(ctx)
			if err != nil {
				log.Printf("Error consuming message: %v", err)
				continue
			}

			// log/show  key  if needed
			if key != nil {
				log.Printf("Processing order with key: %s", string(key))
			}

			var order Order
			if err := json.Unmarshal(value, &order); err != nil {
				log.Printf("Error unmarshaling order: %v", err)
				continue
			}

			// Update analytics
			a.analytics.TotalOrders++
			a.analytics.TotalRevenue += order.Amount
			a.analytics.ProductCounts[order.ProductID]++
			hour := order.Timestamp.Hour()
			a.analytics.HourlyRevenue[hour] += order.Amount

			// Log/show current analytics
			a.logCurrentAnalytics()
		}
	}
}

func (a *Analyzer) logCurrentAnalytics() {
	fmt.Printf("\n=== Current Analytics ===\n")
	fmt.Printf("Total Orders: %d\n", a.analytics.TotalOrders)
	fmt.Printf("Total Revenue: $%.2f\n", a.analytics.TotalRevenue)
	fmt.Printf("Product Counts: %v\n", a.analytics.ProductCounts)
	fmt.Printf("Hourly Revenue: %v\n", a.analytics.HourlyRevenue)
}
