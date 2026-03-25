package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Client es la instancia global del cliente de MongoDB.
// Se inicializa al llamar a Connect y se reutiliza en toda la aplicación.
var Client *mongo.Client

// Connect establece la conexión con MongoDB usando la URI proporcionada.
// Devuelve un error si la conexión falla, o nil si fue exitosa.
func Connect(uri string) error {

	// Configura las opciones del cliente con la URI de conexión (host, puerto, credenciales, etc.)
	clientOptions := options.Client().ApplyURI(uri)

	// Crea el cliente y establece la conexión con MongoDB
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return err
	}

	log.Println("Conectado a MongoDB correctamente")

	// Guarda el cliente en la variable global para usarlo en el resto de la aplicación
	Client = client
	return nil
}

//Función que recibe la colleción completa de la base de datos
func GetCollection(dbName, collectionName string) *mongo.Collection{
	return Client.Database(dbName).Collection(collectionName)
}

func Disconnect() {
	ctx,cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := Client.Disconnect(ctx); err != nil {
		log.Println("Error al desconectar MongoDB:",err)
	}

	log.Println("Desconectado de MongoDB")
}