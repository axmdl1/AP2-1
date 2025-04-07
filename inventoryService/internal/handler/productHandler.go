package handler

import (
	"AP-1/inventoryService/internal/entity"
	"AP-1/inventoryService/internal/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
)

type ProductHandler struct {
	usecase usecase.ProductUsecase
}

func NewProductHandler(u usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{usecase: u}
}

func (h *ProductHandler) StorePage(c *gin.Context) {
	products, err := h.usecase.ListProducts(c.Request.Context(), 0, 100)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not load products"})
		return
	}

	c.HTML(http.StatusOK, "store.html", gin.H{
		"products": products,
	})
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "store.html", nil)
		return
	}

	var product entity.Product
	if err := c.ShouldBind(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.usecase.CreateProduct(c.Request.Context(), &product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, "/products/store")
}

func (h *ProductHandler) GetEditPage(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	product, err := h.usecase.GetProductByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	c.HTML(http.StatusOK, "edit.html", gin.H{"product": product})
}

type UpdateProductForm struct {
	Name        string  `form:"name"`
	Category    string  `form:"category"`
	Description string  `form:"description"`
	Price       float64 `form:"price"`
	Stock       int     `form:"stock"`
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	var form UpdateProductForm

	// Bind form data (excluding the id field)
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the product ID from the form data as a single value.
	idStr := c.PostForm("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product id is required"})
		return
	}
	idStr = strings.TrimSpace(idStr)

	oid, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	// Create the product from form data
	product := entity.Product{
		ID:          oid,
		Name:        form.Name,
		Category:    form.Category,
		Description: form.Description,
		Price:       form.Price,
		Stock:       form.Stock,
	}

	if err := h.usecase.UpdateProduct(c.Request.Context(), &product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/products/store")
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.PostForm("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	id = strings.TrimSpace(id)

	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.usecase.DeleteProduct(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/products/store")
}
