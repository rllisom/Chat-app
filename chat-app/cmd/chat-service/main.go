package chatservice

import (
	"chat-app/internal/chat"
	"chat-app/internal/chatservice"
	"chat-app/internal/db"
	pb "chat-app/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

const mongoURI = "mongo://localhost:27017"

func main(){

	if err := db.Connect(mongoURI); err != nil {
		log.Fatal("Error a conectase con la base de datos ",err.Error())
	}

	defer db.Disconnect()

	chatRepo := chat.NewChatRepository()
	chatService := chat.NewChatService(chatRepo)
	grpcServer := chatservice.NewServer(chatService)
	server := grpc.NewServer()
	pb.RegisterChatServiceServer(server,grpcServer)

	listener, err := net.Listen("tcp",":9002")
	if err != nil {
		log.Fatal("Error abriendo el puerto 9002: ",err)
	}

	log.Println("Chat Service gRPC escuchando en :9002")

	if err := server.Serve(listener); err != nil{
		log.Fatal("Error en gRPC server")
	}
}