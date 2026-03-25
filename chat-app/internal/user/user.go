package user

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string `bson:"username" json:"username"`
	Email string `bson:"email" json:"email"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

func NewUser (username, email string) User {
	return User{
		ID: primitive.NewObjectID(),
		Username: username,
		Email: email,
		CreatedAt: time.Now(),
	}
	
}