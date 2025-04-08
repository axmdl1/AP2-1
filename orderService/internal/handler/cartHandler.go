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

func (h *OrderHandler) ShowPaymentPage(c *gin.Context) {
	var total float64
	for _, item := range cartItems {
		total += item.Price * float64(item.Quantity)
	}
	c.HTML(http.StatusOK, "payment.html", gin.H{"total": total})
}

// ProcessPayment handles POST /orders/checkout
func (h *OrderHandler) ProcessPayment(c *gin.Context) {
	// Bind payment details from the form.
	cardNumber := c.PostForm("card_number")
	expiry := c.PostForm("expiry")
	cvv := c.PostForm("cvv")
	if cardNumber == "" || expiry == "" || cvv == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing payment details"})
		return
	}
	// For demonstration, assume payment processing is successful.

	var total float64
	for _, item := range cartItems {
		total += item.Price * float64(item.Quantity)
	}
	if len(cartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cart is empty"})
		return
	}
	order := entity.Order{
		UserID:     "demo_user", // In a real app, get this from user session/data.
		Products:   cartItems,
		TotalPrice: total,
		Status:     "completed",
		CreatedAt:  time.Now(),
	}

	if err := h.usecase.CreateOrder(c.Request.Context(), &order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Clear the cart.
	cartItems = nil

	// Render a success page.
	c.HTML(http.StatusOK, "buySuccess.html", gin.H{
		"orderID":    order.ID.Hex(),
		"totalPrice": order.TotalPrice,
	})
}
