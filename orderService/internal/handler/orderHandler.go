package handler

import (
	"context"
	"log"

	"AP-1/orderService/internal/entity"
	"AP-1/orderService/internal/usecase"
	pb "AP-1/pb/orderService"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type OrderServiceServer struct {
	pb.UnimplementedOrderServiceServer
	usecase usecase.OrderUsecase
}

func NewOrderServiceServer(u usecase.OrderUsecase) *OrderServiceServer {
	return &OrderServiceServer{usecase: u}
}

func (s *OrderServiceServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	log.Printf("CreateOrder request: %+v", req)
	order := entity.Order{
		UserID: req.Order.UserId,
		Status: req.Order.Status,
	}
	for _, item := range req.Order.Items {
		order.Items = append(order.Items, entity.OrderItem{
			ProductID: item.ProductId,
			Quantity:  int(item.Quantity),
			Price:     item.Price,
		})
	}

	var total float64
	for _, item := range order.Items {
		total += item.Price * float64(item.Quantity)
	}
	order.TotalPrice = total
	order.CreatedAt = time.Now()

	if err := s.usecase.CreateOrder(ctx, &order); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	resOrder := &pb.Order{
		Id:         order.ID,
		UserId:     order.UserID,
		Status:     order.Status,
		TotalPrice: order.TotalPrice,
		CreatedAt:  order.CreatedAt.Format(time.RFC3339),
	}
	for _, item := range order.Items {
		resOrder.Items = append(resOrder.Items, &pb.OrderItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
			Price:     item.Price,
		})
	}

	return &pb.CreateOrderResponse{
		Order:   resOrder,
		Message: "Order created successfully",
	}, nil
}

func (s *OrderServiceServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	log.Printf("GetOrder request: %+v", req)
	order, err := s.usecase.GetOrderByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	resOrder := &pb.Order{
		Id:         order.ID,
		UserId:     order.UserID,
		Status:     order.Status,
		TotalPrice: order.TotalPrice,
		CreatedAt:  order.CreatedAt.Format(time.RFC3339),
	}
	for _, item := range order.Items {
		resOrder.Items = append(resOrder.Items, &pb.OrderItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
			Price:     item.Price,
		})
	}
	return &pb.GetOrderResponse{Order: resOrder}, nil
}

func (s *OrderServiceServer) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	log.Printf("ListOrders request: %+v", req)
	orders, err := s.usecase.ListOrders(ctx, req.Skip, req.Limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	var resOrders []*pb.Order
	for _, order := range orders {
		resOrder := &pb.Order{
			Id:         order.ID,
			UserId:     order.UserID,
			Status:     order.Status,
			TotalPrice: order.TotalPrice,
			CreatedAt:  order.CreatedAt.Format(time.RFC3339),
		}
		for _, item := range order.Items {
			resOrder.Items = append(resOrder.Items, &pb.OrderItem{
				ProductId: item.ProductID,
				Quantity:  int32(item.Quantity),
				Price:     item.Price,
			})
		}
		resOrders = append(resOrders, resOrder)
	}
	return &pb.ListOrdersResponse{Orders: resOrders}, nil
}

func (s *OrderServiceServer) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.UpdateOrderResponse, error) {
	log.Printf("UpdateOrder request: %+v", req)
	order := entity.Order{
		ID:         req.Order.Id,
		UserID:     req.Order.UserId,
		Status:     req.Order.Status,
		TotalPrice: req.Order.TotalPrice,
	}
	for _, item := range req.Order.Items {
		order.Items = append(order.Items, entity.OrderItem{
			ProductID: item.ProductId,
			Quantity:  int(item.Quantity),
			Price:     item.Price,
		})
	}
	var total float64
	for _, item := range order.Items {
		total += item.Price * float64(item.Quantity)
	}
	order.TotalPrice = total

	if err := s.usecase.UpdateOrder(ctx, &order); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	resOrder := &pb.Order{
		Id:         order.ID,
		UserId:     order.UserID,
		Status:     order.Status,
		TotalPrice: order.TotalPrice,
		CreatedAt:  order.CreatedAt.Format(time.RFC3339),
	}
	for _, item := range order.Items {
		resOrder.Items = append(resOrder.Items, &pb.OrderItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
			Price:     item.Price,
		})
	}

	return &pb.UpdateOrderResponse{
		Order:   resOrder,
		Message: "Order updated successfully",
	}, nil
}

func (s *OrderServiceServer) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	log.Printf("DeleteOrder request: %+v", req)
	if err := s.usecase.DeleteOrder(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.DeleteOrderResponse{Message: "Order deleted successfully"}, nil
}
