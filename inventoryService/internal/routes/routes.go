package routes

import (
	"AP-1/inventoryService/internal/handler"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, ph *handler.ProductHandler) {
	router.POST("/products", ph.CreateProduct)
}
