package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/static/", http.FileServer(http.Dir("public")))
	log.Print("Running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
