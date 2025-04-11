package main

import (
	"AP-1/api-gateway/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	//middlewares
	//router.Use(middleware.AuthMiddleware)

	//router
	routes.SetupRoutes(router)

	//server port
	router.Run(":1000")
}
