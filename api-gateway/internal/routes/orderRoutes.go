package routes

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	pb "AP-1/pb/orderService"
	"github.com/gin-gonic/gin"
)

// RegisterOrderRoutes registers REST endpoints for order operations.
func RegisterOrderRoutes(router *gin.Engine, client pb.OrderServiceClient) {
	// Create Order
	router.POST("/orders", func(c *gin.Context) {
		var req pb.CreateOrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		res, err := client.CreateOrder(ctx, &req)
		if err != nil {
			log.Printf("CreateOrder error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, res)
	})

	// Get Order by ID
	router.GET("/orders/:id", func(c *gin.Context) {
		id := c.Param("id")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		res, err := client.GetOrder(ctx, &pb.GetOrderRequest{Id: id})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.HTML(http.StatusOK, "edit_order.html", gin.H{"order": res.Order})
	})

	// List Orders (with pagination)
	router.GET("/orders", func(c *gin.Context) {
		skipStr := c.Query("skip")
		limitStr := c.Query("limit")
		skip, _ := strconv.Atoi(skipStr)
		limit, _ := strconv.Atoi(limitStr)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		req := &pb.ListOrdersRequest{
			Skip:  int32(skip),
			Limit: int32(limit),
		}
		res, err := client.ListOrders(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		//c.JSON(http.StatusOK, res.Orders)
		c.HTML(http.StatusOK, "orders.html", gin.H{"orders": res.Orders})
	})

	// Update Order
	router.POST("/orders/edit", func(c *gin.Context) {
		// Bind only the new status and order id from the form.
		var req struct {
			ID     string `form:"id" binding:"required"`
			Status string `form:"status" binding:"required"`
		}
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Fetch the existing order details using the gRPC GetOrder method.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		existingRes, err := client.GetOrder(ctx, &pb.GetOrderRequest{Id: req.ID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve existing order: " + err.Error()})
			return
		}

		updatedOrder := &pb.Order{
			Id:         existingRes.Order.Id,
			UserId:     existingRes.Order.UserId,
			Items:      existingRes.Order.Items,
			TotalPrice: existingRes.Order.TotalPrice,
			CreatedAt:  existingRes.Order.CreatedAt,
			Status:     req.Status,
		}

		grpcReq := &pb.UpdateOrderRequest{
			Order: updatedOrder,
		}

		ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel2()
		_, err = client.UpdateOrder(ctx2, grpcReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Redirect(http.StatusMovedPermanently, "/orders")
	})

	// Delete Order
	router.POST("/orders/:id", func(c *gin.Context) {
		id := c.Param("id")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := client.DeleteOrder(ctx, &pb.DeleteOrderRequest{Id: id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Redirect(http.StatusMovedPermanently, "/orders")
	})

	// --- Cart and Checkout Endpoints ---

	// Add item to cart
	router.POST("/orders/cart/add", func(c *gin.Context) {
		var req pb.AddToCartRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		res, err := client.AddToCart(ctx, &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	})

	// View cart items
	router.GET("/orders/cart", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		res, err := client.ViewCart(ctx, &pb.ViewCartRequest{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	})

	// Show payment page (in API Gateway we can forward as plain JSON or HTML; here we use JSON for simplicity)
	router.GET("/orders/checkout", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// ViewCartResponse already includes total
		res, err := client.ViewCart(ctx, &pb.ViewCartRequest{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Returning JSON for front-end to render payment page.
		//c.JSON(http.StatusOK, gin.H{"total": res.Total})
		c.HTML(http.StatusOK, "checkout.html", gin.H{"total": res.Total})
	})

	// Process checkout (payment)
	router.POST("/orders/checkout", func(c *gin.Context) {
		var req pb.CheckoutRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		res, err := client.Checkout(ctx, &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	})
}
