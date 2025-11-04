package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var UserCollection *mongo.Collection
var TodoCollection *mongo.Collection
var PreferencesCollection *mongo.Collection

func Connect(uri string, dbName string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	Client = client
	UserCollection = client.Database(dbName).Collection("users")
	TodoCollection = client.Database(dbName).Collection("todos")
	PreferencesCollection = client.Database(dbName).Collection("preferences") // 👈 agregamos esta línea
}

// 🔹 Función auxiliar para obtener cualquier colección por nombre
func GetCollection(name string) *mongo.Collection {
	return Client.Database("todoapp").Collection(name)
}
