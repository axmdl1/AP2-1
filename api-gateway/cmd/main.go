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
	// Dial the User Service (gRPC)
	userConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to User Service: %v", err)
	}
	defer userConn.Close()
	userClient := pbUser.NewUserServiceClient(userConn)

	// Dial the Inventory Service (gRPC)
	invConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to Inventory Service: %v", err)
	}
	defer invConn.Close()
	invClient := pbInventory.NewInventoryServiceClient(invConn)

	// Dial the Order Service (gRPC)
	orderConn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to Order Service: %v", err)
	}
	defer orderConn.Close()
	orderClient := pbOrder.NewOrderServiceClient(orderConn)

	// Initialize the Gin router.
	router := gin.Default()

	// Set up CORS for allowing our frontend (adjust origin as needed)
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:1001"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Serve UI (HTML pages) from the ui folder.
	router.Static("/ui", "./apiGateway/ui")
	// Optionally, if you want to serve a default landing page:
	router.GET("/", func(c *gin.Context) {
		c.File("./apiGateway/ui/index.html")
	})

	// Set up API routes for each microservice.
	// These route functions are defined in separate files under internal/routes.
	routes.RegisterUserRoutes(router, userClient)
	routes.RegisterInventoryRoutes(router, invClient)
	routes.RegisterOrderRoutes(router, orderClient)

	// Start API Gateway on port 1004.
	log.Println("API Gateway starting on port 1004...")
	if err := router.Run(":1004"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
