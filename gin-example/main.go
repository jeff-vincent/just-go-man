package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gin-gonic/gin"
)

func insertDoc(client mongo.Client, key string, value string) {
	coll := client.Database("example").Collection("example")
	_, err := coll.InsertOne(context.TODO(), bson.D{{key, value}})

	if err != nil {
		fmt.Println(err)
	}

}

func getDoc(client mongo.Client, key string, value string) bson.D {
	coll := client.Database("example").Collection("example")
	var result bson.D
	err := coll.FindOne(context.TODO(), bson.D{{key, value}}).Decode(&result)

	if err != nil {
		fmt.Println(err)
	}

	return result

}

func main() {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
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

	r.POST("/get-doc", func(c *gin.Context) {
		key := c.PostForm("key")
		value := c.PostForm("value")
		result := getDoc(*client, key, value)

		c.JSON(http.StatusOK, gin.H{
			"message": result,
		})
	})
	r.Run()

}
