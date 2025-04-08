package handler

import (
	"AP-1/orderService/internal/entity"
	"AP-1/orderService/internal/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
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

func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	var form struct {
		ID     string `form:"id" binding:"required"`
		Status string `form:"status" binding:"required"`
	}

	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	oid, err := primitive.ObjectIDFromHex(form.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	order := entity.Order{
		ID:     oid,
		Status: form.Status,
	}

	if err := h.usecase.UpdateOrder(c.Request.Context(), &order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/orders")
}

// ListOrders handles GET /orders
func (h *OrderHandler) ListOrders(c *gin.Context) {
	orders, err := h.usecase.ListOrders(c.Request.Context(), 0, 100)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not list orders"})
		return
	}
	c.HTML(http.StatusOK, "orders.html", gin.H{"orders": orders})
}

func (h *OrderHandler) GetEditOrderPage(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order id is required"})
		return
	}

	order, err := h.usecase.GetOrderByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.HTML(http.StatusOK, "edit_order.html", gin.H{"order": order})
}

func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	id := c.PostForm("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order id is required"})
		return
	}
	if err := h.usecase.DeleteOrder(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, "/orders")
}
