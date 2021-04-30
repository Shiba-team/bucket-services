package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Client Database instance

var BucketCollection, FileCollection *mongo.Collection

func ConnectDatabase() {
	MongoDb := os.Getenv("MONGODB_URL")
	fmt.Println(MongoDb)
	Client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	err = Client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	BucketCollection = Client.Database("storage").Collection("bucket")
	FileCollection = Client.Database("storage").Collection("file")
}
