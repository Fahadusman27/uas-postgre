package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB *mongo.Database

// ConnectMongoDB membuat koneksi ke MongoDB
func ConnectMongoDB() *mongo.Database {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017" // Default
	}

	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = "achievements_db" // Default
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("MongoDB connection failed: ", err)
	}

	// Ping database untuk memastikan koneksi berhasil
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("MongoDB ping failed: ", err)
	}

	MongoDB = client.Database(dbName)
	log.Println("MongoDB connected successfully to database:", dbName)

	return MongoDB
}

// GetMongoDB mengembalikan instance MongoDB database
func GetMongoDB() *mongo.Database {
	return MongoDB
}
