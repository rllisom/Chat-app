package main

import (
	"chat-app/internal/db"
	"chat-app/internal/user"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	//1.Conectar a MongoDB

	err := db.Connect("mongodb://localhost:27017")

	if err != nil {
		log.Fatal("No se pudo conectar a MongoDB", err)
	}
	defer db.Disconnect()

	userRepo := user.NewUserRepository()
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)

	router := gin.Default()

	v1 := router.Group("/api/v1")
	userHandler.RegisterRoutes(v1)

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Error al arrancar el servidor: ",err)
	}

}