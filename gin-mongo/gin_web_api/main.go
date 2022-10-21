package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type Doc struct {
	Data bson.D
}

type Docs struct {
	Data []bson.M
}

var DB_SERVICE_HOST = os.Getenv("DB_SERVICE_HOST")
var DB_SERVICE_PORT = os.Getenv("DB_SERVICE_PORT")

func insertDoc(c *gin.Context) {
	key := c.PostForm("key")
	value := c.PostForm("value")
	formData := url.Values{"key": {key}, "value": {value}}
	address := fmt.Sprintf("http://%s:%s/insert-doc", DB_SERVICE_HOST, DB_SERVICE_PORT)
	resp, _ := http.PostForm(address, formData)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	c.JSON(http.StatusOK, gin.H{
		"Data": body,
	})
}

func getDoc(c *gin.Context) {
	key := c.Query("key")
	value := c.Query("value")
	address := fmt.Sprintf("http://%s:%s/get-doc?key=%s&value=%s", DB_SERVICE_HOST, DB_SERVICE_PORT, key, value)
	resp, _ := http.Get(address)
	defer resp.Body.Close()
	val := &Doc{}
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(val)
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"Data": val,
	})
}

func getAllDocs(c *gin.Context) {
	address := fmt.Sprintf("http://%s:%s/get-all-docs", DB_SERVICE_HOST, DB_SERVICE_PORT)
	resp, _ := http.Get(address)
	defer resp.Body.Close()
	val := &Docs{}
	fmt.Println(resp.Body)
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(val)
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"Data": val,
	})
}

func main() {
	r := gin.Default()
	r.POST("/insert-doc", insertDoc)
	r.GET("/get-doc", getDoc)
	r.GET("/get-all-docs", getAllDocs)
	r.Run(":8081")
}
