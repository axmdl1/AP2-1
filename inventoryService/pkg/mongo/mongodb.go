package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var client *mongo.Client

/*func Connect() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO-URI")) // Replace with your MongoDB URI

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
		return nil, err
	}

	// Check connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
		return nil, err
	}

	log.Println("Successfully connected to MongoDB")
	return client, nil
}*/

func ConnectMongoDB(uri, dbName string) *mongo.Database {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	return client.Database(dbName)
}
