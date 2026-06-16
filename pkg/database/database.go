package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func init() {
	REQUIRED_ENV := []string{"DATABASE_URL", "DATABASE_NAME"}
	for _, key := range REQUIRED_ENV {
		if key == "" {
			log.Fatalf("missing required env var: %s", key)
		}
	}
}

func CreateClient() (client *mongo.Client, database *mongo.Database, err error) {
	DATABASE_NAME := os.Getenv("DATABASE_NAME")
	DATABASE_URL := os.Getenv("DATABASE_URL")
	_, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	client, err = mongo.Connect(options.Client().ApplyURI(DATABASE_URL))
	database = client.Database(DATABASE_NAME)
	return
}
