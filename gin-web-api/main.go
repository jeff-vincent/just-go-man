package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type Doc struct {
	Data bson.D
}

type Docs struct {
	Data []bson.M
}

func insertDoc(c *gin.Context) {
	key := c.PostForm("key")
	value := c.PostForm("value")
	formData := url.Values{"key": {key}, "value": {value}}
	resp, _ := http.PostForm("http://127.0.0.1:8080/insert-doc", formData)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	c.JSON(http.StatusOK, gin.H{
		"Data": body,
	})
}

func getDoc(c *gin.Context) {
	key := c.Query("key")
	value := c.Query("value")
	address := fmt.Sprintf("http://127.0.0.1:8080/get-doc?key=%s&value=%s", key, value)
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
	resp, _ := http.Get("http://127.0.0.1:8080/get-all-docs")
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
