package chat

import (
	pb "chat-app/proto"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type GRPCHandler struct {
	client pb.ChatServiceClient
}

func NewGRPCHandler(client pb.ChatServiceClient) *GRPCHandler {
	return &GRPCHandler{client: client}
}

func (h *GRPCHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rooms := rg.Group("/rooms")
	{
		rooms.POST("", h.CreateRoom)
		rooms.GET("", h.GetAllRooms)
		rooms.POST("/:id/messages", h.SendMessage)
		rooms.GET("/:id/messages", h.GetRoomMessages)
	}
}

func (h *GRPCHandler) CreateRoom(c *gin.Context) {
	var req struct {
		Name      string `json:"name"       binding:"required"`
		CreatedBy string `json:"created_by" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	room, err := h.client.CreateRoom(ctx, &pb.CreateRoomRequest{
		Name:      req.Name,
		CreatedBy: req.CreatedBy,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, room)
}

func (h *GRPCHandler) GetAllRooms(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.client.GetAllRooms(ctx, &pb.GetAllRoomsRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp.Rooms)
}

func (h *GRPCHandler) SendMessage(c *gin.Context) {
	roomID := c.Param("id")

	var req struct {
		UserID   string `json:"user_id"  binding:"required"`
		Username string `json:"username" binding:"required"`
		Content  string `json:"content"  binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	msg, err := h.client.SendMessage(ctx, &pb.SendMessageRequest{
		RoomId:   roomID,
		UserId:   req.UserID,
		Username: req.Username,
		Content:  req.Content,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, msg)
}

func (h *GRPCHandler) GetRoomMessages(c *gin.Context) {
	roomID := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.client.GetRoomMessages(ctx, &pb.GetRoomMessagesRequest{RoomId: roomID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp.Messages)
}