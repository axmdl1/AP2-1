package handler

import (
	"AP-1/orderService/internal/entity"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

var cartItems []entity.OrderItem

// AddToCart handles POST /orders/cart/add
func (h *OrderHandler) AddToCart(c *gin.Context) {
	productID := c.PostForm("product_id")
	quantityStr := c.PostForm("quantity")
	priceStr := c.PostForm("price")
	if productID == "" || quantityStr == "" || priceStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing product_id, quantity or price"})
		return
	}
	q, err := strconv.Atoi(quantityStr)
	if err != nil || q <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quantity"})
		return
	}
	p, err := strconv.ParseFloat(priceStr, 64)
	if err != nil || p < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid price"})
		return
	}
	orderItem := entity.OrderItem{
		ProductID: productID,
		Quantity:  q,
		Price:     p,
	}
	cartItems = append(cartItems, orderItem)
	c.JSON(http.StatusOK, gin.H{"message": "item added to cart", "cart": cartItems})
}

// ViewCart handles GET /orders/cart
func (h *OrderHandler) ViewCart(c *gin.Context) {
	var total float64
	for _, item := range cartItems {
		total += item.Price * float64(item.Quantity)
	}
	c.HTML(http.StatusOK, "cart.html", gin.H{
		"cartItems": cartItems,
		"total":     total,
	})
}

// BuyCart handles POST /orders/cart/buy
func (h *OrderHandler) BuyCart(c *gin.Context) {
	if len(cartItems) == 0 {
		c.HTML(http.StatusBadRequest, "cart.html", gin.H{
			"error":     "Cart is empty",
			"cartItems": cartItems,
			"total":     0,
		})
		return
	}
	var total float64
	for _, item := range cartItems {
		total += item.Price * float64(item.Quantity)
	}

	order := entity.Order{
		UserID:     "demo_user", // This can be dynamic in a real application.
		Products:   cartItems,
		TotalPrice: total,
		Status:     "pending",
		CreatedAt:  time.Now(),
	}

	if err := h.usecase.CreateOrder(c.Request.Context(), &order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Clear the cart after placing the order.
	cartItems = nil
	c.HTML(http.StatusOK, "buySuccess.html", gin.H{
		"orderID":    order.ID.Hex(),
		"totalPrice": order.TotalPrice,
	})
}
