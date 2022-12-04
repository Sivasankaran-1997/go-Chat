package domain

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DB *mongo.Database
)

func ConnDB() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	uri := os.Getenv("MONGO_DB")

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*20)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	DB = client.Database("chat")
	fmt.Println("Successfuly connected to the database.")
}
