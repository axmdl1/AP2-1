package main

import (
	"AP-1/inventoryService/internal/handler"
	"AP-1/inventoryService/internal/repository"
	"AP-1/inventoryService/internal/routes"
	"AP-1/inventoryService/internal/usecase"
	"AP-1/inventoryService/pkg/mongo"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Using system environment variables.")
	}
	database := mongo.ConnectMongoDB("mongodb://localhost:27017/inventoryDB", "inventoryDB")

	//init
	productRepo := repository.NewProductRepository(database)
	productUsecase := usecase.NewProductUsecase(productRepo)
	productHandler := handler.NewProductHandler(productUsecase)

	//router
	router := gin.Default()
	routes.SetupRoutes(router, productHandler)

	router.Run(":1001")
}
