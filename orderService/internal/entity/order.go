package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" form:"id"`
	UserID     string             `bson:"user_id" json:"user_id" form:"user_id"`
	Products   []OrderItem        `bson:"products" json:"products" form:"products"`
	Status     string             `bson:"status" json:"status" form:"status"`
	TotalPrice float64            `bson:"total_price" json:"total_price" form:"total_price"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at" form:"created_at"`
}

type OrderItem struct {
	ProductID string  `bson:"product_id" json:"product_id" form:"product_id"`
	Quantity  int     `bson:"quantity" json:"quantity" form:"quantity"`
	Price     float64 `bson:"price" json:"price" form:"price"` // Price at the time of order
}
