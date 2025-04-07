package repository

import (
	"AP-1/inventoryService/internal/entity"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ProductRepository defines the CRUD operations for products.
type ProductRepository interface {
	Create(ctx context.Context, product *entity.Product) error
	FindByID(ctx context.Context, id string) (*entity.Product, error)
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter bson.M, skip, limit int64) ([]entity.Product, error)
}

// productRepositoryMongo implements ProductRepository using MongoDB.
type productRepositoryMongo struct {
	collection *mongo.Collection
}

// NewProductRepository creates a new repository instance.
func NewProductRepository(db *mongo.Database) ProductRepository {
	return &productRepositoryMongo{
		collection: db.Collection("products"),
	}
}

func (r *productRepositoryMongo) Create(ctx context.Context, product *entity.Product) error {
	res, err := r.collection.InsertOne(ctx, product)
	if err != nil {
		return err
	}
	product.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *productRepositoryMongo) FindByID(ctx context.Context, id string) (*entity.Product, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var product entity.Product
	if err := r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&product); err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepositoryMongo) Update(ctx context.Context, product *entity.Product) error {
	filter := bson.M{"_id": product.ID}
	update := bson.M{"$set": product}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *productRepositoryMongo) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

func (r *productRepositoryMongo) List(ctx context.Context, filter bson.M, skip, limit int64) ([]entity.Product, error) {
	opts := options.Find().SetSkip(skip).SetLimit(limit)
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
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return products, nil
}
