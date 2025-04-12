package repository

import (
	"AP-1/inventoryService/internal/entity"
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository interface {
	Create(ctx context.Context, product *entity.Product) error
	FindByID(ctx context.Context, id string) (*entity.Product, error)
	List(ctx context.Context, filter bson.M, skip, limit int32) ([]entity.Product, error)
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, id string) error
}

type productRepositoryMongo struct {
	collection *mongo.Collection
}

func NewProductRepository(db *mongo.Database) ProductRepository {
	return &productRepositoryMongo{
		collection: db.Collection("products"),
	}
}

func (r *productRepositoryMongo) Create(ctx context.Context, product *entity.Product) error {
	_, err := r.collection.InsertOne(ctx, product)
	return err
}

func (r *productRepositoryMongo) FindByID(ctx context.Context, id string) (*entity.Product, error) {
	var product entity.Product
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepositoryMongo) List(ctx context.Context, filter bson.M, skip, limit int32) ([]entity.Product, error) {
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	cursor, err := r.collection.Find(ctx, filter, opts)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var products []entity.Product
	for cursor.Next(ctx) {
		var product entity.Product
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (r *productRepositoryMongo) Update(ctx context.Context, product *entity.Product) error {
	filter := bson.M{"id": product.ID}
	update := bson.M{"$set": bson.M{
		"name":        product.Name,
		"category":    product.Category,
		"description": product.Description,
		"price":       product.Price,
		"stock":       product.Stock,
	}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *productRepositoryMongo) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}

func int64Ptr(i int64) *int64 {
	return &i
}
