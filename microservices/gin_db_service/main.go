package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gin-gonic/gin"
)

var MONGO_HOST = os.Getenv("MONGO_HOST")
var MONGO_PORT = os.Getenv("MONGO_PORT")

func insertDoc(client mongo.Client, key string, value string) {
	coll := client.Database("example").Collection("example")
	_, err := coll.InsertOne(context.TODO(), bson.D{{key, value}})

	if err != nil {
		log.Fatal(err)
	}
}

func getDoc(client mongo.Client, key string, value string) bson.D {
	coll := client.Database("example").Collection("example")
	var result bson.D
	err := coll.FindOne(context.TODO(), bson.D{{key, value}}).Decode(&result)

	if err != nil {
		log.Fatal(err)
	}

	return result
}

func getAllDocs(client mongo.Client) []bson.M {
	coll := client.Database("example").Collection("example")
	cursor, err := coll.Find(context.TODO(), bson.D{})

	if err != nil {
		log.Fatal(err)
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	return results
}

func main() {
	mongo_uri := fmt.Sprintf("mongodb://%s:%s", MONGO_HOST, MONGO_PORT)
	client, err := mongo.NewClient(options.Client().ApplyURI(mongo_uri))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	r := gin.Default()

	r.POST("/insert-doc", func(c *gin.Context) {
		key := c.PostForm("key")
		value := c.PostForm("value")
		insertDoc(*client, key, value)
	})

	r.GET("/get-doc", func(c *gin.Context) {
		key := c.Query("key")
		value := c.Query("value")
		result := getDoc(*client, key, value)

		c.JSON(http.StatusOK, gin.H{
			"Data": result,
		})
	})

	r.GET("/get-all-docs", func(c *gin.Context) {
		result := getAllDocs(*client)
		c.JSON(http.StatusOK, gin.H{
			"Data": result,
		})
	})
	r.Run()
}
