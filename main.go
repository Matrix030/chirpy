package main

import (
	"log"
	"net/http"
)

func main() {
	serveMux := http.NewServeMux()
	server := &http.Server{}
	server.Addr = ":8080"
	server.Handler = serveMux
	serveMux.Handle("/", http.FileServer(http.Dir(".")))
	log.Fatal(server.ListenAndServe())

}
