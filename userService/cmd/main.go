package main

import (
	"AP-1/userService/internal/handler"
	"AP-1/userService/internal/repository"
	"AP-1/userService/internal/routes"
	"AP-1/userService/internal/usecase"
	"AP-1/userService/pkg/mongo"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	db := mongo.ConnectMongoDB("mongodb://localhost:27017/userServiceDB", "userServiceDB")

	userRepo := repository.NewUserRepository(db)
	useCase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(useCase)

	//router
	router := gin.Default()
	routes.SetupRoutes(router, userHandler)

	router.LoadHTMLGlob("userService/ui/*")
	router.Static("/ui", "./ui")

	if err := router.Run(":1004"); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
