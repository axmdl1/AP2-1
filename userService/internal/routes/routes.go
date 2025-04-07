package routes

import (
	"AP-1/userService/internal/handler"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, handler *handler.UserHandler) {
	users := router.Group("/users")
	{
		users.GET("/register", handler.Register)
		users.POST("/register", handler.Register)
		users.GET("/login", handler.Login)
		users.POST("/login", handler.Login)
	}
}
