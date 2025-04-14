package usecase

import (
	"AP-1/orderService/internal/entity"
	"AP-1/orderService/internal/repository"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"time"

	"github.com/google/uuid"
)

type OrderUsecase interface {
	CreateOrder(ctx context.Context, order *entity.Order) error
	GetOrderByID(ctx context.Context, id string) (*entity.Order, error)
	ListOrders(ctx context.Context, skip, limit int32) ([]entity.Order, error)
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
	if len(order.Items) == 0 {
		return errors.New("order must have at least one item")
	}
	var total float64
	for _, item := range order.Items {
		total += item.Price * float64(item.Quantity)
	}
	order.TotalPrice = total

	order.ID = uuid.New().String()
	if order.Status == "" {
		order.Status = "pending"
	}
	order.CreatedAt = time.Now()

	return u.repo.Create(ctx, order)
}

func (u *orderUsecase) GetOrderByID(ctx context.Context, id string) (*entity.Order, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *orderUsecase) ListOrders(ctx context.Context, skip, limit int32) ([]entity.Order, error) {
	filter := bson.M{}
	return u.repo.List(ctx, filter, skip, limit)
}

func (u *orderUsecase) UpdateOrder(ctx context.Context, order *entity.Order) error {
	return u.repo.Update(ctx, order)
}

func (u *orderUsecase) DeleteOrder(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}
