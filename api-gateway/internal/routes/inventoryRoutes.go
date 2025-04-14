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
		// Parse form-encoded data
		var req struct {
			Name        string  `form:"name" binding:"required"`
			Category    string  `form:"category" binding:"required"`
			Description string  `form:"description" binding:"required"`
			Price       float64 `form:"price" binding:"required"`
			Stock       int     `form:"stock" binding:"required"`
		}
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Build the gRPC request.
		grpcReq := &pb.CreateProductRequest{
			Product: &pb.Product{
				Name:        req.Name,
				Category:    req.Category,
				Description: req.Description,
				Price:       req.Price,
				Stock:       int32(req.Stock),
			},
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := client.CreateProduct(ctx, grpcReq)
		if err != nil {
			log.Printf("CreateProduct error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Redirect(http.StatusMovedPermanently, "/products")
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
		//c.JSON(http.StatusOK, res)
		c.HTML(http.StatusOK, "edit.html", gin.H{"product": res.Product})
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
		//c.JSON(http.StatusOK, gin.H{"products": res.Products})
		c.HTML(http.StatusOK, "store.html", gin.H{"products": res.Products})
	})

	router.POST("/products/edit", func(c *gin.Context) {
		var req struct {
			ID          string  `form:"id" binding:"required"`
			Name        string  `form:"name" binding:"required"`
			Category    string  `form:"category" binding:"required"`
			Description string  `form:"description" binding:"required"`
			Price       float64 `form:"price" binding:"required"`
			Stock       int     `form:"stock" binding:"required"`
		}
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		grpcReq := &pb.UpdateProductRequest{
			Product: &pb.Product{
				Id:          req.ID,
				Name:        req.Name,
				Category:    req.Category,
				Description: req.Description,
				Price:       req.Price,
				Stock:       int32(req.Stock),
			},
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := client.UpdateProduct(ctx, grpcReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Redirect(http.StatusMovedPermanently, "/products")
	})

	// Delete Product
	router.POST("/products/:id", func(c *gin.Context) {
		id := c.Param("id")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := client.DeleteProduct(ctx, &pb.DeleteProductRequest{Id: id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Redirect(http.StatusMovedPermanently, "/products")
	})
}
