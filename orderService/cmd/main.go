package main

import (
	"AP-1/orderService/internal/handler"
	"AP-1/orderService/internal/repository"
	"AP-1/orderService/internal/routes"
	"AP-1/orderService/internal/usecase"
	"AP-1/orderService/pkg/mongo"
	"github.com/gin-contrib/cors"
	"html/template"
	"log"
	"time"

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

	router.SetFuncMap(template.FuncMap{
		"mul": func(a float64, b int) float64 {
			return a * float64(b)
		},
	})

	routes.SetupRoutes(router, orderHandler)

	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:1001/products/store"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(config))

	router.LoadHTMLGlob("orderService/ui/*")
	router.Static("/ui", "./ui")

	if err := router.Run(":1003"); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
