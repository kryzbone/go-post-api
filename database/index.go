package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"example.com/go-post-api/helpers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx context.Context
var cancel context.CancelFunc
var client *mongo.Client
var err error

func ConnectDB() (*mongo.Database, context.Context) {
	var mongoUri = os.Getenv("MONGODB_URI")
	if mongoUri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable.")
	}
	// create context
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	// ctx = context.Background()

	//Connect to mongoDB
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))
	helpers.CheckErr(err)

	database := client.Database(helpers.DB_NAME)
	return database, ctx
}

func CloseDB() {
	if client == nil {
		return
	}
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		log.Fatal(err)
	}

	// TODO optional you can log your closed MongoDB client
	fmt.Println("Connection to MongoDB closed.")
}
