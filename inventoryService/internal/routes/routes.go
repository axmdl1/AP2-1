package routes

import (
	"AP-1/inventoryService/internal/handler"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, ph *handler.ProductHandler) {
	products := router.Group("/products")
	{
		products.GET("/store", ph.StorePage)
		products.POST("/create", ph.CreateProduct)
		products.POST("/delete", ph.DeleteProduct)
	}
}
