package infrastructor

import (
	"bookify/internal/config"
	"context"
	"fmt"
	mongodriven "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

func NewMongoDatabase(env *config.Database) *mongodriven.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	var mongodbURI string
	if env.DBUser != "" && env.DBPassword != "" {
		mongodbURI = fmt.Sprintf("mongodb+srv://%s:%s@andrew.8ulkv.mongodb.net/?retryWrites=true&w=majority", env.DBUser, env.DBPassword)
	} else {
		mongodbURI = fmt.Sprintf("mongodb://%s:%s", env.DBHost, env.DBPort)
	}

	mongoCon := options.Client().ApplyURI(mongodbURI).SetServerAPIOptions(serverAPI)
	client, err := mongodriven.Connect(ctx, mongoCon)
	if err != nil {
		log.Fatal("error while connecting with mongo:", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("error while trying to ping mongo:", err)
	}

	log.Println("Connected to MongoDB!")
	return client
}

func CloseMongoDBConnection(client *mongodriven.Client) {
	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connection to MongoDB closed.")
}
