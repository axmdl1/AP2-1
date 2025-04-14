package main

import (
	"log"
	"time"

	pbInventory "AP-1/pb/inventoryService"
	pbOrder "AP-1/pb/orderService"
	pbUser "AP-1/pb/userService"

	"AP-1/api-gateway/internal/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	// Dial gRPC services.
	userConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to User Service: %v", err)
	}
	defer userConn.Close()
	userClient := pbUser.NewUserServiceClient(userConn)

	invConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to Inventory Service: %v", err)
	}
	defer invConn.Close()
	invClient := pbInventory.NewInventoryServiceClient(invConn)

	orderConn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to Order Service: %v", err)
	}
	defer orderConn.Close()
	orderClient := pbOrder.NewOrderServiceClient(orderConn)

	// Initialize Gin router.
	router := gin.Default()

	// Set up CORS (allow requests from your frontend origin).
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:1001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Load HTML templates from ui directory.
	// Adjust the path based on your directory structure.
	router.LoadHTMLGlob("api-gateway/ui/*")

	// Register REST endpoints for each microservice.
	// (Assume you have similar route registration files for users, inventory, orders.)
	routes.RegisterUserRoutes(router, userClient)
	routes.RegisterInventoryRoutes(router, invClient)
	routes.RegisterOrderRoutes(router, orderClient)

	// Register frontend UI routes.
	routes.RegisterUIRoutes(router)

	log.Println("API Gateway starting on port 1004...")
	if err := router.Run(":1004"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
