package main

import (
	"log"
	"net"

	"AP-1/orderService/internal/handler"
	"AP-1/orderService/internal/repository"
	"AP-1/orderService/internal/usecase"
	pb "AP-1/pb/orderService"
	"AP-1/userService/pkg/mongo"

	"google.golang.org/grpc"
)

func main() {
	db := mongo.ConnectMongoDB("mongodb://localhost:27017/orderServiceDB", "orderServiceDB")

	orderRepo := repository.NewOrderRepository(db)
	orderUsecase := usecase.NewOrderUsecase(orderRepo)
	orderServiceServer := handler.NewOrderServiceServer(orderUsecase)

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, orderServiceServer)
	log.Println("Order Service gRPC server started on port :50053")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
