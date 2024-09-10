package main

import (
	"log"
	"net/http"
)

type server int

// func (s server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	log.Println("Request received", r.URL.Path)
// 	w.Write([]byte("Hello, world!"))
// }

func main() {
	// var s server
	
	http.ListenAndServe("localhost:8080", &s)
}
