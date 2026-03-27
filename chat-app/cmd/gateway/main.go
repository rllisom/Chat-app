package main

import (
	"chat-app/internal/chat"
	"chat-app/internal/chatservice"
	"chat-app/internal/user"
	"chat-app/internal/userservice"
	"log"

	"github.com/gin-gonic/gin"

)

func main() {
	
	userClient,err := userservice.NewClient("localhost:9001") 

	if err != nil {
		log.Fatal("User Service: ",err)
		
	}
	defer userClient.Close()

	chatClient, err := chatservice.NewClient("localhost:9002")

	if err != nil {
		log.Fatal("Chat Service: ",err)
	}

	defer chatClient.Close()

	userHandler := user.NewGRPCHandler(userClient.GRPCClient())
	chatHandler := chat.NewGRPCHandler(chatClient.GRPCClient())

	router := gin.Default()
	v1 := router.Group("/api/v1")
	userHandler.RegisterRoutes(v1)
	chatHandler.RegisterRoutes(v1)

	log.Println("Gateway en http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Gateway: ", err)
	}
}