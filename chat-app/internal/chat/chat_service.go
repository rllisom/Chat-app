package chat

import (
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatService struct {
	chatRepo *ChatRepository
}

func NewChatService(chatRepo *ChatRepository) *ChatService {
	return &ChatService{chatRepo: chatRepo}
}

// CreateRoom valida y crea una sala nueva.
func (s *ChatService) CreateRoom(name, createdByID string) (*Room, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return nil, errors.New("el nombre de la sala no puede estar vacío")
	}
	if len(name) < 3 {
		return nil, errors.New("el nombre debe tener al menos 3 caracteres")
	}

	userID, err := primitive.ObjectIDFromHex(createdByID)
	if err != nil {
		return nil, errors.New("ID de usuario no válido")
	}

	room := NewRoom(name, userID)
	return s.chatRepo.CreateRoom(room)
}

// GetAllRooms devuelve todas las salas.
func (s *ChatService) GetAllRooms() (*[]Room, error) {
	return s.chatRepo.FindAllRooms()
}

// SendMessage valida y guarda un mensaje.
func (s *ChatService) SendMessage(roomID, userID, username, content string) (*Message, error) {
	content = strings.TrimSpace(content)

	if content == "" {
		return nil, errors.New("el mensaje no puede estar vacío")
	}
	if len(content) > 500 {
		return nil, errors.New("el mensaje no puede superar 500 caracteres")
	}

	roomObjID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return nil, errors.New("ID de sala no válido")
	}

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("ID de usuario no válido")
	}

	msg := NewMessage(roomObjID, userObjID, username, content)
	return s.chatRepo.SaveMessage(msg)
}

// GetRoomMessages devuelve los últimos mensajes de una sala.
func (s *ChatService) GetRoomMessages(roomID string) ([]Message, error) {
	if roomID == "" {
		return nil, errors.New("el ID de sala no puede estar vacío")
	}
	// Últimos 50 mensajes — valor habitual en chats
	return s.chatRepo.FindMessagesByRoom(roomID, 50)
}