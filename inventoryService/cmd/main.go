package main

import (
	"log"
	"net"

	"AP-1/inventoryService/internal/handler"
	"AP-1/inventoryService/internal/repository"
	"AP-1/inventoryService/internal/usecase"
	"AP-1/inventoryService/pkg/mongo"
	pb "AP-1/pb/inventoryService"

	"google.golang.org/grpc"
)

func main() {
	db := mongo.ConnectMongoDB("mongodb://localhost:27017/inventoryDB", "inventoryDB")

	productRepo := repository.NewProductRepository(db)
	productUsecase := usecase.NewProductUsecase(productRepo)
	inventoryServiceServer := handler.NewInventoryServiceServer(productUsecase)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterInventoryServiceServer(grpcServer, inventoryServiceServer)
	log.Println("Inventory Service gRPC server started on port :50052")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
