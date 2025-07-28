package main

import (
	"log"
	"net/http"
)

func main() {
	Servemux := http.NewServeMux()
	server := &http.Server{}
	server.Addr = ":8080"
	server.Handler = Servemux
	log.Fatal(http.ListenAndServe(":8080", nil))
}
