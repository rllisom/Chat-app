package user

import (
	pb "chat-app/proto"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type GRPCHandler struct {
	client pb.UserServiceClient
}

func NewGRPCHandler(client pb.UserServiceClient) *GRPCHandler {
	return &GRPCHandler{client: client}
}

func (h *GRPCHandler) RegisterRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	{
		users.POST("", h.CreateUser)
		users.GET("", h.GetAllUsers)
		users.GET("/:id", h.GetUserByID)
	}
}

func (h *GRPCHandler) CreateUser(c *gin.Context) {
	var req pb.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := h.client.CreateUser(ctx, &pb.CreateUserRequest{
		Username: req.Username,
		Email:    req.Email,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (h *GRPCHandler) GetAllUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.client.GetAllUsers(ctx, &pb.GetAllUsersRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp.Users)
}

func (h *GRPCHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := h.client.GetUserById(ctx, &pb.GetUserByIdRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}