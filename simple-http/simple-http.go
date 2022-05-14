package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func getAlbums(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(albums)
	if err != nil {
		fmt.Println("Unable to convert the struct to a JSON string")
	} else {
		// convert []byte to a string type and then print
		fmt.Fprint(w, string(b))
	}

}

func GetAlbumByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/albums/")
	// for loop looking for matching ID
	for _, a := range albums {
		if a.ID == id {
			// convert matching struct to JSON
			b, err := json.Marshal(a)
			if err != nil {
				fmt.Println("whoops...")
			} else {
				//write response
				fmt.Fprint(w, string(b))
			}
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func stuff(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "this is different...")
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/different", stuff)
	http.HandleFunc("/albums", getAlbums)
	http.HandleFunc("/albums/", GetAlbumByID)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
