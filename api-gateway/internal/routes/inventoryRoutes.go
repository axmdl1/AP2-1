package routes

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	pb "AP-1/pb/inventoryService"
	"github.com/gin-gonic/gin"
)

// RegisterInventoryRoutes registers REST endpoints for inventory operations.
func RegisterInventoryRoutes(router *gin.Engine, client pb.InventoryServiceClient) {
	// Create a Product
	router.POST("/products", func(c *gin.Context) {
		var req pb.CreateProductRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		res, err := client.CreateProduct(ctx, &req)
		if err != nil {
			log.Printf("CreateProduct error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, res)
	})

	// Get a Product by ID
	router.GET("/products/:id", func(c *gin.Context) {
		id := c.Param("id")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		res, err := client.GetProduct(ctx, &pb.GetProductRequest{Id: id})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	})

	// List Products with pagination
	router.GET("/products", func(c *gin.Context) {
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
		res, err := client.ListProducts(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	})

	// Update Product
	router.PUT("/products/:id", func(c *gin.Context) {
		id := c.Param("id")
		var req pb.UpdateProductRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Ensure the product id in request matches the URL.
		req.Product.Id = id
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		res, err := client.UpdateProduct(ctx, &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	})

	// Delete Product
	router.DELETE("/products/:id", func(c *gin.Context) {
		id := c.Param("id")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		res, err := client.DeleteProduct(ctx, &pb.DeleteProductRequest{Id: id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	})
}
