package usecase

import (
	"AP-1/inventoryService/internal/entity"
	"AP-1/inventoryService/internal/repository"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
)

type ProductUsecase interface {
	CreateProduct(ctx context.Context, product *entity.Product) error
	ListProducts(ctx context.Context, skip, limit int64) ([]*entity.Product, error)
	DeleteProduct(ctx context.Context, id string) error
}

type productUsecase struct {
	repo repository.ProductRepository
}

func NewProductUsecase(repo repository.ProductRepository) ProductUsecase {
	return &productUsecase{repo: repo}
}

func (u productUsecase) CreateProduct(ctx context.Context, product *entity.Product) error {
	if product.Price <= 0 {
		return errors.New("price must be positive")
	}

	if product.Stock < 0 {
		return errors.New("stock cannot be negative")
	}

	return u.repo.Create(ctx, product)
}

func (u *productUsecase) ListProducts(ctx context.Context, skip, limit int64) ([]*entity.Product, error) {
	filter := bson.M{}
	products, err := u.repo.List(ctx, filter, skip, limit)
	if err != nil {
		return nil, err
	}

	// Convert []entity.Product to []*entity.Product
	ptrProducts := make([]*entity.Product, len(products))
	for i := range products {
		ptrProducts[i] = &products[i]
	}
	return ptrProducts, nil
}

func (u *productUsecase) DeleteProduct(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}
