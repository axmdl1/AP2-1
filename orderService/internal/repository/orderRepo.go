package repository

import (
	"AP-1/orderService/internal/entity"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderRepository interface {
	Create(ctx context.Context, order *entity.Order) error
	FindByID(ctx context.Context, id string) (*entity.Order, error)
	Update(ctx context.Context, order *entity.Order) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter bson.M, skip, limit int32) ([]entity.Order, error)
}

type orderRepositoryMongo struct {
	collection *mongo.Collection
}

func NewOrderRepository(db *mongo.Database) OrderRepository {
	return &orderRepositoryMongo{
		collection: db.Collection("orders"),
	}
}

func (r *orderRepositoryMongo) Create(ctx context.Context, order *entity.Order) error {
	_, err := r.collection.InsertOne(ctx, order)
	return err
}

func (r *orderRepositoryMongo) FindByID(ctx context.Context, id string) (*entity.Order, error) {
	var order entity.Order
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepositoryMongo) Update(ctx context.Context, order *entity.Order) error {
	filter := bson.M{"id": order.ID}
	update := bson.M{"$set": bson.M{
		"user_id":     order.UserID,
		"items":       order.Items,
		"total_price": order.TotalPrice,
		"status":      order.Status,
	}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *orderRepositoryMongo) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}

func (r *orderRepositoryMongo) List(ctx context.Context, filter bson.M, skip, limit int32) ([]entity.Order, error) {
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var orders []entity.Order
	for cursor.Next(ctx) {
		var order entity.Order
		if err := cursor.Decode(&order); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}
