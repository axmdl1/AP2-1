package routes

import (
	"AP-1/orderService/internal/handler"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, oh *handler.OrderHandler) {
	orders := router.Group("/orders")
	{
		orders.POST("/create", oh.CreateOrder)
		orders.GET("/:id", oh.GetOrder)
		orders.GET("", oh.ListOrders)
		orders.GET("/edit", oh.GetEditOrderPage)
		orders.POST("/edit", oh.UpdateOrder)
		orders.POST("/delete", oh.DeleteOrder)

		orders.POST("/cart/add", oh.AddToCart)
		orders.GET("/cart", oh.ViewCart)

		orders.GET("/checkout", oh.ShowPaymentPage)
		orders.POST("/checkout", oh.ProcessPayment)
	}
}
