package mongo

import (
	"context"
	"log"
	"time"

	"cinema-booking/pkg/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func Connect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.App.MongoURI))
	if err != nil {
		log.Fatal("Mongo connection error:", err)
	}

	DB = client.Database(config.App.MongoDB)
	log.Println("MongoDB connected")
}
