package usecase

import (
	"AP-1/inventoryService/internal/entity"
	"AP-1/inventoryService/internal/repository"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/google/uuid"
)

type ProductUsecase interface {
	CreateProduct(ctx context.Context, product *entity.Product) error
	GetProductByID(ctx context.Context, id string) (*entity.Product, error)
	ListProducts(ctx context.Context, skip, limit int32) ([]entity.Product, error)
	UpdateProduct(ctx context.Context, product *entity.Product) error
	DeleteProduct(ctx context.Context, id string) error
}

type productUsecase struct {
	repo repository.ProductRepository
}

func NewProductUsecase(repo repository.ProductRepository) ProductUsecase {
	return &productUsecase{repo: repo}
}

func (u *productUsecase) CreateProduct(ctx context.Context, product *entity.Product) error {
	if product.Price <= 0 {
		return errors.New("price must be positive")
	}
	if product.Stock < 0 {
		return errors.New("stock cannot be negative")
	}
	product.ID = uuid.New().String()
	return u.repo.Create(ctx, product)
}

func (u *productUsecase) GetProductByID(ctx context.Context, id string) (*entity.Product, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *productUsecase) ListProducts(ctx context.Context, skip, limit int32) ([]entity.Product, error) {
	filter := bson.M{}
	return u.repo.List(ctx, filter, skip, limit)
}

func (u *productUsecase) UpdateProduct(ctx context.Context, product *entity.Product) error {
	return u.repo.Update(ctx, product)
}

func (u *productUsecase) DeleteProduct(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}
