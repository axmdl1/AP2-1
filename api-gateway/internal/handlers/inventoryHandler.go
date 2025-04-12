package handler

import (
	"context"
	"log"
	"strconv"
	"time"

	pb "AP-1/pb/inventoryService"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	client pb.InventoryServiceClient
}

func NewInventoryHandler(client pb.InventoryServiceClient) *InventoryHandler {
	return &InventoryHandler{client: client}
}

func (h *InventoryHandler) CreateProduct(c *gin.Context) {
	var req pb.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.client.CreateProduct(ctx, &req)
	if err != nil {
		log.Printf("CreateProduct error: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, res)
}

func (h *InventoryHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.client.GetProduct(ctx, &pb.GetProductRequest{Id: id})
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

func (h *InventoryHandler) ListProducts(c *gin.Context) {
	skipStr := c.Query("skip")
	limitStr := c.Query("limit")
	skip, _ := strconv.Atoi(skipStr)
	limit, _ := strconv.Atoi(limitStr)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.ListProductsRequest{
		Skip:  int32(skip),
		Limit: int32(limit),
	}
	res, err := h.client.ListProducts(ctx, req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

func (h *InventoryHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var req pb.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	req.Product.Id = id // Ensure the product ID matches the URL.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.client.UpdateProduct(ctx, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

func (h *InventoryHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.client.DeleteProduct(ctx, &pb.DeleteProductRequest{Id: id})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}
