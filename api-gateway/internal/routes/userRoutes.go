package routes

import (
	"context"
	"log"
	"time"

	pb "AP-1/pb/userService"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes registers REST endpoints for user operations.
func RegisterUserRoutes(router *gin.Engine, client pb.UserServiceClient) {
	router.POST("/users/register", func(c *gin.Context) {
		var req pb.RegisterUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		res, err := client.RegisterUser(ctx, &req)
		if err != nil {
			log.Printf("RegisterUser error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, res)
	})

	router.POST("/users/login", func(c *gin.Context) {
		var req pb.AuthenticateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		res, err := client.AuthenticateUser(ctx, &req)
		if err != nil {
			log.Printf("AuthenticateUser error: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	})

	router.GET("/users/profile", func(c *gin.Context) {
		userID := c.Query("user_id")
		if userID == "" || userID == "undefined" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		res, err := client.GetUserProfile(ctx, &pb.GetUserProfileRequest{UserID: userID})
		if err != nil {
			log.Printf("GetUserProfile error: %v", err)
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	})
}
