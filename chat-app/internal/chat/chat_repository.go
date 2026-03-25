package chat

import (
	"chat-app/internal/db"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ChatRepository struct {
	rooms *mongo.Collection
	messages *mongo.Collection
}

func NewChatRepository() *ChatRepository {
	return &ChatRepository{
		rooms: db.GetCollection("chat_app","rooms"),
		messages: db.GetCollection("chat_app","messages"),
	}
}

//------------------Rooms-----------------------

func (r *ChatRepository) CreateRoom(room Room) (*Room,error){
	ctx,cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	if _,err := r.rooms.InsertOne(ctx,room); err != nil {
		return nil, errors.New("Error al crear un nuevo chat")
	}

	return &room,nil
}

func (r *ChatRepository) FindRoomByID(id string) (*Room,error){
	ctx,cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	objectID,err  := primitive.ObjectIDFromHex(id)

	if err != nil{
		return nil, errors.New("ID del chat no válido")
	}
	var room Room

	if err := r.rooms.FindOne(ctx,objectID).Decode(&room); err != nil {
		if errors.Is(err,mongo.ErrNoDocuments){
			return nil, errors.New("Chat no encontrado")
		}
		return nil,err
	}

	return &room,nil
}

func (r *ChatRepository) FindAllRooms() (*[]Room,error){

	ctx,cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	var rooms []Room

	cursor, err := r.rooms.Find(ctx,bson.M{})

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err:= cursor.All(ctx,&rooms); err != nil {
		return nil,err
	}
	return &rooms,nil
}

//--------------Messages------------------------------------
func(r *ChatRepository) SaveMessage(m Message)(*Message,error){
	ctx,cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	_,err := r.messages.InsertOne(ctx,m)

	if err != nil {
		return nil,err
	}

	return &m,nil
}

func(r *ChatRepository)FindMessagesByRoom(roomID string, limit int) ([]Message, error){
	ctx,cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	objectID,err := primitive.ObjectIDFromHex(roomID)

	if err != nil {
		return nil,err
	}

	opt := options.Find().
	SetSort(bson.D{{Key:"sent_at", Value: -1}}).
	SetLimit(int64(limit))

	cursor,err := r.messages.Find(ctx,bson.M{"room_id":objectID},opt)
	if err != nil {
		return nil,err
	}

	var messages []Message

	if err := cursor.All(ctx,&messages); err != nil {
		return nil,err
	}

	return messages,nil

}