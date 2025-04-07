package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id" form:"id"`
	Name        string             `form:"name" json:"name"`
	Category    string             `form:"category" json:"category"`
	Description string             `form:"description" json:"description"`
	Price       float64            `form:"price" json:"price"`
	Stock       int                `form:"stock" json:"stock"`
}
