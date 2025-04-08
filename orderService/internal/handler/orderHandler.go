package handler

import (
	"AP-1/orderService/internal/entity"
	"AP-1/orderService/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderHandler struct {
	usecase usecase.OrderUsecase
}

func NewOrderHandler(u usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{usecase: u}
}

// CreateOrder handles POST /orders/create
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var order entity.Order
	if err := c.ShouldBind(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Optionally, add server-side validation for order fields

	if err := h.usecase.CreateOrder(c.Request.Context(), &order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, order)
}

// GetOrder handles GET /orders/:id
func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	order, err := h.usecase.GetOrderByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

// UpdateOrder handles POST /orders/edit
func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	var order entity.Order
	if err := c.ShouldBind(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idStr := c.PostForm("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order id is required"})
		return
	}
	oid, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}
	order.ID = oid

	if err := h.usecase.UpdateOrder(c.Request.Context(), &order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order updated"})
}

// ListOrders handles GET /orders
func (h *OrderHandler) ListOrders(c *gin.Context) {
	orders, err := h.usecase.ListOrders(c.Request.Context(), 0, 100)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not list orders"})
		return
	}
	c.JSON(http.StatusOK, orders)
}
