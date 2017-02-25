package main

import (
	"log"
	"net/http"
)

func main() {
	PORT := "127.0.0.1:8000"
	log.Print("Running server on " + PORT)

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/static/", fs)

	http.HandleFunc("/", BoringHandler)
	log.Fatal(http.ListenAndServe(PORT, nil))
}

// BoringHandler tells user to go to cat
func BoringHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is boring go to /static/img/cat.jpg"))
}
