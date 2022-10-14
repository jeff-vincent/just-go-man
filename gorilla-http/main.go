package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type album struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{}

func getAlbums(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(albums)
	if err != nil {
		fmt.Println("Unable to convert the struct to a JSON string")
	} else {
		fmt.Fprint(w, string(b))
	}
}

func getAlbumByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/albums/")
	id_as_int, _ := strconv.Atoi(id)
	for _, a := range albums {
		if a.ID == id_as_int {
			b, err := json.Marshal(&a)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Fprint(w, string(b))
			}
		}
	}
}

func createAlbum(latest_id int, w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	artist := r.FormValue("artist")
	p := r.FormValue("price")
	price, _ := strconv.ParseFloat(p, 64)
	id := latest_id + 1
	a := album{ID: id, Title: title, Artist: artist, Price: price}
	albums = append(albums, a)
	b, _ := json.Marshal(a)
	fmt.Fprint(w, string(b))
}

func saveAlbum(w http.ResponseWriter, r *http.Request) {
	if len(albums) < 1 {
		latest_id := 0
		createAlbum(latest_id, w, r)
	} else {
		latest_id := albums[len(albums)-1].ID
		createAlbum(latest_id, w, r)
	}
}

func remove(a album, albs *[]album) []album {
	for index, alb := range albums {
		if alb == a {
			arr := make([]album, 0)
			arr = append(arr, albums[0:index]...)
			albums = append(arr, albums[index+1:]...)

			return albums
		}
	}
	return nil
}

func deleteAlbum(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/delete-album/")
	id_as_int, _ := strconv.Atoi(id)
	for _, a := range albums {
		if a.ID == id_as_int {
			albs := &albums
			albums := remove(a, albs)
			b, err := json.Marshal(&albums)
			if err != nil {
				return
			}
			fmt.Fprint(w, string(b))
		}
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}
func main() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/albums", getAlbums)
	http.HandleFunc("/albums/", getAlbumByID)
	http.HandleFunc("/save-album", saveAlbum)
	http.HandleFunc("/delete-album/", deleteAlbum)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
