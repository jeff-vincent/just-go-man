package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
)

type BlogPost struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Body   string `json:"body"`
}

type Doc struct {
	Data bson.D `json:"data"`
}

type Docs struct {
	Data []bson.M `json:"data"`
}

var REDIS_HOST = os.Getenv("REDIS_HOST")
var REDIS_PORT = os.Getenv("REDIS_PORT")
var SERVICE_HOST = os.Getenv("SERVICE_HOST")
var SERVICE_PORT = os.Getenv("SERVICE_PORT")

func getDoc(c *gin.Context) {
	title := c.Query("title")
	address := fmt.Sprintf("http://%s:%s/get-doc?title=%s", SERVICE_HOST, SERVICE_PORT, title)
	resp, _ := http.Get(address)
	defer resp.Body.Close()
	val := &Doc{}
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(val)

	if err != nil {
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, val)
}

func getAllDocs(c *gin.Context) {
	address := fmt.Sprintf("http://%s:%s/get-all-docs", SERVICE_HOST, SERVICE_PORT)
	resp, _ := http.Get(address)
	defer resp.Body.Close()
	val := &Docs{}
	fmt.Println(resp.Body)
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(val)
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, val)
}

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func newPost(c *gin.Context, t string, a string, b string) {
	c.HTML(http.StatusOK, "post.html", gin.H{
		"title":  t,
		"author": a,
		"body":   b,
	})
}

func main() {
	redis_uri := fmt.Sprintf("redis://%s:%s/0", REDIS_HOST, REDIS_PORT)
	opt, err := redis.ParseURL(redis_uri)
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(opt)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")

	r.GET("/", index)
	r.POST("/insert-doc", func(c *gin.Context) {
		title := c.PostForm("title")
		author := c.PostForm("author")
		body := c.PostForm("body")
		new_post := BlogPost{Title: title, Author: author, Body: body}
		payload, err := json.Marshal(new_post)
		if err != nil {
			fmt.Println(err)
		}
		ctx := context.Background()
		err = rdb.Publish(ctx, "1", payload).Err()
		if err != nil {
			panic(err)
		}
		newPost(c, title, author, body)
	})
	r.GET("/get-doc", getDoc)
	r.GET("/get-all-docs", getAllDocs)
	r.Run(":8081")
}
