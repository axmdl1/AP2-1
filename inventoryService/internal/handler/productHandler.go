package handler

import (
	"AP-1/inventoryService/internal/entity"
	"AP-1/inventoryService/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ProductHandler struct {
	usecase *usecase.ProductUsecase
}

func NewProductHandler(u *usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{usecase: u}
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
