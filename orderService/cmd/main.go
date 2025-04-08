package main

import (
	"AP-1/orderService/internal/handler"
	"AP-1/orderService/internal/repository"
	"AP-1/orderService/internal/routes"
	"AP-1/orderService/internal/usecase"
	"AP-1/orderService/pkg/mongo"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"time"
)

func main() {
	// Connect to MongoDB
	db := mongo.ConnectMongoDB("mongodb://localhost:27017/orderServiceDB", "orderServiceDB")

	// Initialize repository, usecase, and handler
	orderRepo := repository.NewOrderRepository(db)
	orderUsecase := usecase.NewOrderUsecase(orderRepo)
	orderHandler := handler.NewOrderHandler(orderUsecase)

	router := gin.Default()

	// Merge custom template functions into a single FuncMap
	router.SetFuncMap(template.FuncMap{
		"mul": func(a float64, b int) float64 {
			return a * float64(b)
		},
		"selected": func(current, option string) string {
			if current == option {
				return "selected"
			}
			return ""
		},
	})

	// Load templates (adjust path as needed)
	router.LoadHTMLGlob("orderService/ui/*")
	router.Static("/ui", "./ui")

	routes.SetupRoutes(router, orderHandler)

	// CORS configuration
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:1001"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(config))

	if err := router.Run(":1003"); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
