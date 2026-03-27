package chatservice

import (
	"chat-app/internal/chat"
	pb "chat-app/proto"
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


type Server struct {
	pb.UnimplementedChatServiceServer
	service *chat.ChatService
}

func NewServer(service *chat.ChatService) *Server {
	return &Server{service: service}
}


func (s *Server) CreateRoom(ctx context.Context, req *pb.CreateRoomRequest)(*pb.RoomResponse,error){
	room,err := s.service.CreateRoom(req.Name,req.CreatedBy)

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
	}

	return roomToProto(room),nil
}

func (s *Server) GetAllRooms(ctx context.Context, req *pb.GetAllRoomsRequest) (*pb.GetAllRoomsResponse, error) {
	rooms, err := s.service.GetAllRooms()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	var protoRooms []*pb.RoomResponse
	for _, r := range *rooms {
		protoRooms = append(protoRooms, roomToProto(&r))
	}

	return &pb.GetAllRoomsResponse{Rooms: protoRooms}, nil
}

func (s *Server) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.MessageResponse, error) {
	msg, err := s.service.SendMessage(
		req.RoomId,
		req.UserId,
		req.Username,
		req.Content,
	)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
	}

	return messageToProto(msg), nil
}


func (s *Server) GetRoomMessages(ctx context.Context, req *pb.GetRoomMessagesRequest) (*pb.GetRoomMessagesResponse, error) {
	messages, err := s.service.GetRoomMessages(req.RoomId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	var protoMessages []*pb.MessageResponse
	for _, m := range messages {
		protoMessages = append(protoMessages, messageToProto(&m))
	}

	return &pb.GetRoomMessagesResponse{Messages: protoMessages}, nil
}

//-----Helpers-------------------------------

func roomToProto(room *chat.Room) *pb.RoomResponse {
	return &pb.RoomResponse{
		Id: room.ID.Hex(),
		Name: room.Name,
		CreatedBy: room.CreatedBy.Hex(),
		CreatedAt :  room.CreatedAt.Format(time.RFC3339),
	}
}

func messageToProto(m *chat.Message) *pb.MessageResponse{
	return &pb.MessageResponse{
		Id: m.ID.Hex(),
		RoomId: m.RoomID.Hex(),
		UserId: m.UserID.Hex(),
		Content: m.Content,
		SentAt: m.SentAt.Format(time.RFC3339),
	}
}