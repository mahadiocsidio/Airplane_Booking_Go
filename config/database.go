package config

import(
	"context"
	// "os"
	"log"
	"time"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB(connectionString string) *mongo.Client{
	// bikin context dengan timeout biar gak ngegantung
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal("Error connect ke MongoDB:", err)
	}

	// cek apakah koneksi bener2 jalan
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Error ping MongoDB:", err)
	}

	fmt.Println("Connected to MongoDB ✅")
	return client
}


func GetCollection(client *mongo.Client, db string, collectionName string) *mongo.Collection{
	collection := client.Database(db).Collection(collectionName)
	return collection
}