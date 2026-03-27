package userservice

import (
	"chat-app/internal/user"
	pb "chat-app/proto"
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	service *user.UserService
}

func NewServer(service *user.UserService) *Server {
	return &Server{service: service}
}

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse,error){
	u,err := s.service.CreateUser(req.Username,req.Email)

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
	}

	return toProto(u), nil
}

func (s *Server) GetUserByID (ctx context.Context, req *pb.GetUserByIdRequest) (*pb.UserResponse,error){
	u,err := s.service.GetUserByID(req.Id)

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%s", err.Error())
	}

	return toProto(u),nil
}

func (s *Server) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse,error){

	users, err := s.service.GetAllUsers()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	var protoUsers []*pb.UserResponse
	for _, u := range *users {
		protoUsers = append(protoUsers,toProto(&u))
	}

	return &pb.GetAllUsersResponse{Users: protoUsers},nil
}


func toProto(u *user.User) *pb.UserResponse {
	return &pb.UserResponse{
		Id: u.ID.Hex(),
		Username: u.Username,
		Email: u.Email,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
	}
}


