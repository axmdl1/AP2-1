package main

import (
	"AP-1/inventoryService/internal/handler"
	"AP-1/inventoryService/internal/repository"
	"AP-1/inventoryService/internal/routes"
	"AP-1/inventoryService/internal/usecase"
	"AP-1/inventoryService/pkg/mongo"
	"github.com/gin-gonic/gin"
)

func main() {
	database := mongo.ConnectMongoDB("mongodb://localhost:27017/inventoryDB", "inventoryDB")

	//init
	productRepo := repository.NewProductRepository(database)
	productUsecase := usecase.NewProductUsecase(productRepo)
	productHandler := handler.NewProductHandler(productUsecase)

	//router
	router := gin.Default()
	routes.SetupRoutes(router, productHandler)

	router.LoadHTMLGlob("inventoryService/ui/*")
	router.Static("/ui", "./ui")

	router.Run(":1001")
}
