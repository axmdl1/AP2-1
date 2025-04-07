package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" || !strings.HasPrefix(token, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.Next()
}
