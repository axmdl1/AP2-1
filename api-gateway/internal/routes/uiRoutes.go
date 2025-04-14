package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterUIRoutes registers endpoints that serve HTML pages.
func RegisterUIRoutes(router *gin.Engine) {
	// Serve the landing page.
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Registration page.
	router.GET("/users/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	})

	// Login page.
	router.GET("/users/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	// User profile page.
	router.GET("/profile", func(c *gin.Context) {
		// Optionally, you can pass profile data if needed.
		c.HTML(http.StatusOK, "profile.html", nil)
	})

	// Store page for products.
	router.GET("/store", func(c *gin.Context) {
		c.HTML(http.StatusOK, "store.html", nil)
	})

	// Cart page.
	router.GET("/cart", func(c *gin.Context) {
		c.HTML(http.StatusOK, "cart.html", nil)
	})

	// Orders history page.
	/*router.GET("/orders", func(c *gin.Context) {
		c.HTML(http.StatusOK, "orders.html", nil)
	})*/
}
