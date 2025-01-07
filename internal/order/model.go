package order

import "time"

type Order struct {
	OrderID   string    `json:"order_id"`
	UserID    string    `json:"user_id"`
	ProductID string    `json:"product_id"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

type Analytics struct {
	TotalOrders   int             `json:"total_orders"`
	TotalRevenue  float64         `json:"total_revenue"`
	ProductCounts map[string]int  `json:"product_counts"`
	HourlyRevenue map[int]float64 `json:"hourly_revenue"`
}

const (
	StatusPending   = "pending"
	StatusPaid      = "paid"
	StatusShipped   = "shipped"
	StatusDelivered = "delivered"
)
