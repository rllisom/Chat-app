package chatservice

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "chat-app/proto"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.ChatServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:   conn,
		client: pb.NewChatServiceClient(conn),
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) GRPCClient() pb.ChatServiceClient {
	return c.client
}

func (c *Client) CreateRoom(name, createdBy string) (*pb.RoomResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.client.CreateRoom(ctx, &pb.CreateRoomRequest{
		Name:      name,
		CreatedBy: createdBy,
	})
}

func (c *Client) GetAllRooms() (*pb.GetAllRoomsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.client.GetAllRooms(ctx, &pb.GetAllRoomsRequest{})
}

func (c *Client) SendMessage(roomID, userID, username, content string) (*pb.MessageResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.client.SendMessage(ctx, &pb.SendMessageRequest{
		RoomId:   roomID,
		UserId:   userID,
		Username: username,
		Content:  content,
	})
}

func (c *Client) GetRoomMessages(roomID string) (*pb.GetRoomMessagesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.client.GetRoomMessages(ctx, &pb.GetRoomMessagesRequest{
		RoomId: roomID,
	})
}