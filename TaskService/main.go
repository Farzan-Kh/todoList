package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Task struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title    string             `json:"title" bson:"title"`
	Status   string             `json:"status" bson:"status"`
	Username string             `json:"username" bson:"username"`
}

var db *mongo.Collection
var ctx = context.Background()

func initDB() {
	mongoURI := os.Getenv("MONGO_URI")
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic("failed to create MongoDB client")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		panic("failed to connect to MongoDB")
	}

	db = client.Database("taskservice").Collection("tasks")
}

func createTask(c *gin.Context) {
	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.ID = primitive.NewObjectID()
	_, err := db.InsertOne(ctx, task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}
	c.JSON(http.StatusCreated, task)
}

func listTasks(c *gin.Context) {
	username := c.Query("username")
	var tasks []Task
	cursor, err := db.Find(ctx, bson.M{"username": username})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &tasks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse tasks"})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func main() {
	initDB()

	r := gin.Default()
	r.POST("/tasks/", createTask)
	r.GET("/tasks/", listTasks)

	r.Run(":8082") // Run on port 8082
}
