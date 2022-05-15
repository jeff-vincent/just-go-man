package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice literal
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

func getAlbumByID(w http.ResponseWriter, r *http.Request) {
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

func saveAlbum(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	artist := r.FormValue("artist")
	p := r.FormValue("price")
	price, _ := strconv.ParseFloat(p, 64)
	latest_id := albums[len(albums)-1].ID
	id_as_int, _ := strconv.Atoi(latest_id)
	//TODO: fix this
	id := string(id_as_int + 1)
	a := album{ID: id, Title: title, Artist: artist, Price: price}
	albums = append(albums, a)
	b, _ := json.Marshal(a)
	fmt.Fprint(w, string(b))
}

func remove(a album, albums []album) []album {
	for index, alb := range albums {
		if alb == a {
			// TODO: why duplicates of ID 3?
			return append(albums[0:index], albums[index+1:]...)
		}
	}
	return nil
}

func deleteAlbum(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/delete-album/")
	for _, a := range albums {
		if a.ID == id {
			albums := remove(a, albums)
			b, err := json.Marshal(albums)
			if err != nil {
				return
			}
			fmt.Fprint(w, string(b))
		}
	}
}

func updateAlbum(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/update-album/")
	for _, alb := range albums {
		if id == alb.ID {
			p := &alb
			if r.FormValue("artist") != "" {
				p.Artist = r.FormValue("artist")
			}
			if r.FormValue("title") != "" {
				p.Title = r.FormValue("title")
			}
			if r.FormValue("price") != "" {
				p.Price, _ = strconv.ParseFloat(r.FormValue("price"), 64)
			}
		}
		b, _ := json.Marshal(alb)
		fmt.Fprint(w, string(b))
	}

}

func main() {
	http.HandleFunc("/albums", getAlbums)
	http.HandleFunc("/albums/", getAlbumByID)
	http.HandleFunc("/save-album", saveAlbum)
	http.HandleFunc("/delete-album/", deleteAlbum)
	http.HandleFunc("/update-album/", updateAlbum)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
