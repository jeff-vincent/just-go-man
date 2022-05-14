package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func stuff(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "this is different...")
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/different", stuff)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
