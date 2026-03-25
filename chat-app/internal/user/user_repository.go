package user

import (
	"chat-app/internal/db"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository() *UserRepository{
	return &UserRepository{
		collection: db.GetCollection("chat_app","users"),
	}
}

func (r *UserRepository) Create(user User) (*User,error) {
	ctx, cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()
	if _,err := r.collection.InsertOne(ctx,user); err != nil{
		return nil,err
	}
	return &user,nil
}

func (r *UserRepository) FindByID(id string) (*User,error){

	ctx,cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	objectID,err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("ID de usuario no válido")
	}

	var user User

	if err := r.collection.FindOne(ctx,bson.M{"_id":objectID}).Decode(&user); err != nil{
		if errors.Is(err, mongo.ErrNoDocuments){
			return nil, errors.New("Usuario no encontrado")
		}
		return nil,err
	}

	return &user,err
	
}

func (r *UserRepository) FindByUsername(username string) (*User,error){
	ctx,cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	var user User

	if err := r.collection.FindOne(ctx,bson.M{"username":username}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments){
			return nil, errors.New("Usuario no encontrado")
		}

		return nil,err
	}

	return &user,nil

}

func (r *UserRepository) FindAll() (*[]User,error){
	ctx,cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	cursor,err := r.collection.Find(ctx,bson.M{}) //Trae todo de la base de datos de User

	/*
	Usamos cursor para recorrer todos los elementos user que hay en mongo db sin cargarlos en memoria.
	*/
	if err != nil {
		return nil,err
	}
	defer cursor.Close(ctx)

	var users []User

	//Una vez cargado el cursor hacemos el slice a la lista
	if err := cursor.All(ctx,&users); err != nil {
		return nil,err
	}

	return &users,nil
}