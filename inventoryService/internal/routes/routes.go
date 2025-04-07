package routes

import (
	"AP-1/inventoryService/internal/handler"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, ph *handler.ProductHandler) {
	r := router.Group("/products")
	{
		r.GET("/store", ph.CreateProduct)
		r.POST("/create", ph.CreateProduct)
	}
}
