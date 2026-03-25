package chat

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Room struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name string `bson:"name" json:"name"`
	CreatedBy primitive.ObjectID `bson:"created_by" json:"created_by"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RoomID    primitive.ObjectID `bson:"room_id"       json:"room_id"`  
	UserID    primitive.ObjectID `bson:"user_id"       json:"user_id"` 
	Username  string             `bson:"username"      json:"username"` 
	Content   string             `bson:"content"       json:"content"`
	SentAt    time.Time          `bson:"sent_at"       json:"sent_at"`
}

func NewRoom(name string, createdBy primitive.ObjectID) Room {
	return Room{
		ID:        primitive.NewObjectID(),
		Name:      name,
		CreatedBy: createdBy,
		CreatedAt: time.Now(),
	}
}


func NewMessage(roomID, userID primitive.ObjectID, username, content string) Message {
	return Message{
		ID:       primitive.NewObjectID(),
		RoomID:   roomID,
		UserID:   userID,
		Username: username,
		Content:  content,
		SentAt:   time.Now(),
	}
}