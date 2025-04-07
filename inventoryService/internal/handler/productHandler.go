package handler

import (
	"AP-1/inventoryService/internal/entity"
	"AP-1/inventoryService/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
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
		c.HTML(http.StatusOK, "products.html", nil)
		return
	}

	var product entity.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.usecase.CreateProduct(c.Request.Context(), &product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, product)
}
