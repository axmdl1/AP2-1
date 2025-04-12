package handler

import (
	"context"
	"log"
	"strconv"
	"time"

	pb "AP-1/pb/orderService"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	client pb.OrderServiceClient
}

func NewOrderHandler(client pb.OrderServiceClient) *OrderHandler {
	return &OrderHandler{client: client}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req pb.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.client.CreateOrder(ctx, &req)
	if err != nil {
		log.Printf("CreateOrder error: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, res)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.client.GetOrder(ctx, &pb.GetOrderRequest{Id: id})
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

func (h *OrderHandler) ListOrders(c *gin.Context) {
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
	res, err := h.client.ListOrders(ctx, req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	id := c.Param("id")
	var req pb.UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	req.Order.Id = id
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.client.UpdateOrder(ctx, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.client.DeleteOrder(ctx, &pb.DeleteOrderRequest{Id: id})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

//Cart endpoints

func (h *OrderHandler) AddToCart(c *gin.Context) {
	var req pb.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.client.AddToCart(ctx, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

func (h *OrderHandler) ViewCart(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.client.ViewCart(ctx, &pb.ViewCartRequest{})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

func (h *OrderHandler) Checkout(c *gin.Context) {
	var req pb.CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.client.Checkout(ctx, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}
