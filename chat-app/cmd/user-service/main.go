package userservice

import (
	"chat-app/internal/db"
	"chat-app/internal/user"
	"chat-app/internal/userservice"
	pb "chat-app/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

const mongoURI = "mongodb://localhost:27017"

func main() {

	if err := db.Connect(mongoURI); err != nil {
		log.Fatal("Error al conectar MongoDB: ",err)
	}

	defer db.Disconnect()

	userRepo := user.NewUserRepository()
	userService := user.NewUserService(userRepo)
	grpcServer := userservice.NewServer(userService)

	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server,grpcServer)

	listener,err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatal("Error al abrir el puerto 9001")
	}

	log.Println("User Service gRPC escuchando en :9001")
	
	if err := server.Serve(listener); err != nil {
		log.Fatal("Error en gRPC server: ", err)
	}
}