package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoDBClient     *mongo.Client
	MongoDBConnCancel context.CancelFunc

	MongoTheVardiacDB *mongo.Database
)

type MongoDB struct{}

func (cfg *MongoDB) MongoDBMakeConn() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_CONNECTION_URI")))
	if err != nil {
		log.Fatal("Failed to connect to mongo DB: " + err.Error())
	}

	MongoDBClient = client
	MongoDBConnCancel = cancel

	// The vardiac db instance
	MongoTheVardiacDB = client.Database("thevardiac")
}
