package userservice

import (
	pb "chat-app/proto"
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn *grpc.ClientConn
	client pb.UserServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn,err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil,err
	}

	return &Client{
		conn: conn,
		client: pb.NewUserServiceClient(conn),
	},nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) GRPCClient() pb.UserServiceClient {
	return c.client
}


func (c *Client) CreateUser(username,email string) (*pb.UserResponse, error){
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.client.CreateUser(ctx, &pb.CreateUserRequest{
		Username: username,
		Email:    email,
	})
}

func (c *Client) GetUserById(id string) (*pb.UserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.client.GetUserById(ctx, &pb.GetUserByIdRequest{Id: id})
}

func (c *Client) GetAllUsers() (*pb.GetAllUsersResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.client.GetAllUsers(ctx, &pb.GetAllUsersRequest{})
}