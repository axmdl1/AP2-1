package usecase

import (
	"AP-1/inventoryService/internal/entity"
	"AP-1/inventoryService/internal/repository"
	"context"
	"errors"
)

type ProductUsecase struct {
	repo repository.ProductRepository
}

func NewProductUsecase(repo repository.ProductRepository) *ProductUsecase {
	return &ProductUsecase{repo: repo}
}

func (u *ProductUsecase) CreateProduct(ctx context.Context, product *entity.Product) error {
	if product.Price <= 0 {
		return errors.New("price must be positive")
	}

	if product.Stock < 0 {
		return errors.New("stock cannot be negative")
	}

	return u.repo.Create(ctx, product)
}
