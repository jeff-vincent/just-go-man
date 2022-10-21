package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BlogPost struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Body   string `json:"body"`
}

var MONGO1_HOST = os.Getenv("MONGO1_HOST")
var MONGO1_PORT = os.Getenv("MONGO1_PORT")
var MONGO2_HOST = os.Getenv("MONGO2_HOST")
var MONGO2_PORT = os.Getenv("MONGO2_PORT")
var REDIS_HOST = os.Getenv("REDIS_HOST")
var REDIS_PORT = os.Getenv("REDIS_PORT")

func insertDoc(post bson.M, mongo2 *mongo.Client) *mongo.InsertOneResult {
	coll := mongo2.Database("example").Collection("example")
	result, err := coll.InsertOne(context.TODO(), post)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func getDoc(post BlogPost, mongo1 *mongo.Client) bson.M {
	coll := mongo1.Database("example").Collection("example")
	var result []bson.M
	opts := options.Find().SetSort(bson.D{{"age", 1}})
	cursor, err := coll.Find(context.TODO(), bson.D{{"title", post.Title}}, opts)
	if err = cursor.All(context.TODO(), &result); err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	return result[len(result)-1]
}

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	mongo1_uri := fmt.Sprintf("mongodb://%s:%s", MONGO1_HOST, MONGO1_PORT)
	mongo1, err := mongo.NewClient(options.Client().ApplyURI(mongo1_uri))
	if err != nil {
		log.Fatal(err)
	}
	err = mongo1.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer mongo1.Disconnect(ctx)

	mongo2_uri := fmt.Sprintf("mongodb://%s:%s", MONGO2_HOST, MONGO2_PORT)
	mongo2, err := mongo.NewClient(options.Client().ApplyURI(mongo2_uri))
	if err != nil {
		log.Fatal(err)
	}
	err = mongo2.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer mongo2.Disconnect(ctx)

	redis_uri := fmt.Sprintf("redis://%s:%s/0", REDIS_HOST, REDIS_PORT)
	opt, err := redis.ParseURL(redis_uri)
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(opt)
	pubsub := rdb.Subscribe(ctx, "2")
	defer pubsub.Close()
	ch := pubsub.Channel()
	for msg := range ch {
		post := BlogPost{}
		if err := json.Unmarshal([]byte(msg.Payload), &post); err != nil {
			panic(err)
		}
		re := getDoc(post, mongo1)
		insertDoc(re, mongo2)
	}
}
