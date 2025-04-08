package main

import (
	"AP-1/orderService/internal/handler"
	"AP-1/orderService/internal/repository"
	"AP-1/orderService/internal/routes"
	"AP-1/orderService/internal/usecase"
	"AP-1/orderService/pkg/mongo"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to MongoDB
	db := mongo.ConnectMongoDB("mongodb://localhost:27017/orderServiceDB", "orderServiceDB")

	// init
	orderRepo := repository.NewOrderRepository(db)
	orderUsecase := usecase.NewOrderUsecase(orderRepo)
	orderHandler := handler.NewOrderHandler(orderUsecase)

	router := gin.Default()

	routes.SetupRoutes(router, orderHandler)

	log.Println("Starting Order Service on port 1003...")
	if err := router.Run(":1003"); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
