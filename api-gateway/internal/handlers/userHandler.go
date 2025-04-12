package handler

import (
	"context"
	"log"
	"time"

	pb "AP-1/pb/userService"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	client pb.UserServiceClient
}

func NewUserHandler(client pb.UserServiceClient) *UserHandler {
	return &UserHandler{client: client}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req pb.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.client.RegisterUser(ctx, &req)
	if err != nil {
		log.Printf("RegisterUser error: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, res)
}

func (h *UserHandler) AuthenticateUser(c *gin.Context) {
	var req pb.AuthenticateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.client.AuthenticateUser(ctx, &req)
	if err != nil {
		log.Printf("AuthenticateUser error: %v", err)
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

func (h *UserHandler) GetUserProfile(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(400, gin.H{"error": "user_id is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.client.GetUserProfile(ctx, &pb.GetUserProfileRequest{UserID: userID})
	if err != nil {
		log.Printf("GetUserProfile error: %v", err)
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}
