package repository

import (
	"AP-1/orderService/internal/entity"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// OrderRepository defines the operations for orders.
type OrderRepository interface {
	Create(ctx context.Context, order *entity.Order) error
	FindByID(ctx context.Context, id string) (*entity.Order, error)
	Update(ctx context.Context, order *entity.Order) error
	List(ctx context.Context, filter bson.M, skip, limit int64) ([]entity.Order, error)
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
	order.CreatedAt = entity.Order{}.CreatedAt // alternatively, time.Now()
	res, err := r.collection.InsertOne(ctx, order)
	if err != nil {
		return err
	}
	order.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *orderRepositoryMongo) FindByID(ctx context.Context, id string) (*entity.Order, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var order entity.Order
	if err := r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&order); err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepositoryMongo) Update(ctx context.Context, order *entity.Order) error {
	filter := bson.M{"_id": order.ID}
	update := bson.M{"$set": order}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *orderRepositoryMongo) List(ctx context.Context, filter bson.M, skip, limit int64) ([]entity.Order, error) {
	opts := options.Find().SetSkip(skip).SetLimit(limit)
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
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}
