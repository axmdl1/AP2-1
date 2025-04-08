package usecase

import (
	"AP-1/orderService/internal/entity"
	"AP-1/orderService/internal/repository"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
)

type OrderUsecase interface {
	CreateOrder(ctx context.Context, order *entity.Order) error
	GetOrderByID(ctx context.Context, id string) (*entity.Order, error)
	ListOrders(ctx context.Context, skip, limit int64) ([]*entity.Order, error)
	UpdateOrder(ctx context.Context, order *entity.Order) error
	DeleteOrder(ctx context.Context, id string) error
}

type orderUsecase struct {
	repo repository.OrderRepository
}

func NewOrderUsecase(repo repository.OrderRepository) OrderUsecase {
	return &orderUsecase{repo: repo}
}

func (u *orderUsecase) CreateOrder(ctx context.Context, order *entity.Order) error {
	// Basic validation can be added here.
	if order.TotalPrice < 0 {
		return errors.New("total price cannot be negative")
	}
	return u.repo.Create(ctx, order)
}

func (u *orderUsecase) GetOrderByID(ctx context.Context, id string) (*entity.Order, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *orderUsecase) ListOrders(ctx context.Context, skip, limit int64) ([]*entity.Order, error) {
	filter := bson.M{}
	orders, err := u.repo.List(ctx, filter, skip, limit)
	if err != nil {
		return nil, err
	}
	ptrOrders := make([]*entity.Order, len(orders))
	for i := range orders {
		ptrOrders[i] = &orders[i]
	}
	return ptrOrders, nil
}

func (u *orderUsecase) UpdateOrder(ctx context.Context, order *entity.Order) error {
	return u.repo.Update(ctx, order)
}

func (u *orderUsecase) DeleteOrder(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}
