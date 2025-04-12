package handler

import (
	"AP-1/orderService/internal/entity"
	pb "AP-1/pb/orderService"
	"context"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var cartItems []entity.OrderItem

func (s *OrderServiceServer) AddToCart(ctx context.Context, req *pb.AddToCartRequest) (*pb.AddToCartResponse, error) {
	log.Printf("AddToCart request received: %+v", req)
	item := entity.OrderItem{
		ProductID: req.Item.ProductId,
		Quantity:  int(req.Item.Quantity),
		Price:     req.Item.Price,
	}
	cartItems = append(cartItems, item)
	return &pb.AddToCartResponse{
		Message: "Item added to cart",
	}, nil
}

func (s *OrderServiceServer) ViewCart(ctx context.Context, req *pb.ViewCartRequest) (*pb.ViewCartResponse, error) {
	log.Println("ViewCart request received")
	var total float64
	var items []*pb.CartItem
	for _, item := range cartItems {
		total += item.Price * float64(item.Quantity)
		items = append(items, &pb.CartItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
			Price:     item.Price,
		})
	}
	return &pb.ViewCartResponse{
		Items: items,
		Total: total,
	}, nil
}

func (s *OrderServiceServer) Checkout(ctx context.Context, req *pb.CheckoutRequest) (*pb.CheckoutResponse, error) {
	log.Printf("Checkout request received: %+v", req)
	if req.CardNumber == "" || req.Expiry == "" || req.Cvv == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing payment details")
	}

	if len(cartItems) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "cart is empty")
	}

	var total float64
	for _, item := range cartItems {
		total += item.Price * float64(item.Quantity)
	}

	order := entity.Order{
		UserID:     "demo_user",
		Items:      cartItems,
		TotalPrice: total,
		Status:     "completed",
		CreatedAt:  time.Now(),
	}

	if err := s.usecase.CreateOrder(ctx, &order); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	cartItems = nil

	return &pb.CheckoutResponse{
		OrderId:    order.ID,
		TotalPrice: total,
		Message:    "Order placed successfully",
	}, nil
}
