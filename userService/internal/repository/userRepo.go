package repository

import (
	"AP-1/userService/internal/entity"
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByID(ctx context.Context, id string) (*entity.User, error)
}
type userRepositoryMongo struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepositoryMongo{
		collection: db.Collection("users"),
	}
}

func (u userRepositoryMongo) Create(ctx context.Context, user *entity.User) error {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	_, err := u.collection.InsertOne(ctx, user)
	return err
}

func (u userRepositoryMongo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := u.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u userRepositoryMongo) FindByID(ctx context.Context, id string) (*entity.User, error) {
	var user entity.User
	err := u.collection.FindOne(ctx, bson.M{"id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
