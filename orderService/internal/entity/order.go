package entity

import "time"

type OrderItem struct {
	ProductID string  `bson:"product_id" json:"product_id"`
	Quantity  int     `bson:"quantity" json:"quantity"`
	Price     float64 `bson:"price" json:"price"`
}

type Order struct {
	ID         string      `bson:"id" json:"id"`
	UserID     string      `bson:"user_id" json:"user_id"`
	Items      []OrderItem `bson:"items" json:"items"`
	TotalPrice float64     `bson:"total_price" json:"total_price"`
	Status     string      `bson:"status" json:"status"`
	CreatedAt  time.Time   `bson:"created_at" json:"created_at"`
}
