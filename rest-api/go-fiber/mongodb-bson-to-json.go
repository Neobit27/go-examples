package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Task struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	UserId string             `json:"userId"`
	Task   string             `json:"task"`
	Status string             `json:"status"`
}

func main() {
	app := fiber.New()

	clientOptions := options.Client().ApplyURI("mongodburi")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("MongoDB bağlantısı başarılı!")

	taskcollection := client.Database("todolist").Collection("tasks")

	app.Get("/api/tasks", func(c *fiber.Ctx) error {
		var tasks []Task
		cursor, err := taskcollection.Find(context.TODO(), bson.M{})
		if err != nil {
			log.Fatalf("Error finding tasks: %v", err)
		}
		for cursor.Next(context.TODO()) {
			var task Task
			if err := cursor.Decode(&task); err != nil {
				log.Fatalf("Error decoding task: %v", err)
			}
			tasks = append(tasks, task)
		}
		if err := cursor.Err(); err != nil {
			log.Fatalf("Error iterating cursor: %v", err)
		}
		return c.JSON(tasks)
	})

	app.Listen(5000)
}
