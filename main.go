package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

//Create a struct in main.go that will hold any stateful, in-memory data we'll need to keep track of.
//  In our case, we just need to keep track of the number of requests we've received.
type apiConfig struct {
	fileserverHits atomic.Int32
}
func main() {
	const filepathRoot = "."
	const port = "8080"
	apiCfg := apiConfig{}


	mux := http.NewServeMux()
	/*Wrap the http.FileServer handler with the middleware method we just wrote.*/
	fileServer := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fileServer))
	mux.HandleFunc("/metrics", apiCfg.handleMetrics)
	mux.HandleFunc("/reset", apiCfg.handlerReset)
	mux.HandleFunc("/healthz", handlerReadiness)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
/* Next, write a new middleware method on a *apiConfig that increments the fileserverHits counter 
 every time it's called.Here's the method signature I used:*/
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

/*Create a new handler that writes the number of requests that have been counted as plain text in this 
format to the HTTP response:*/
func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	count := cfg.fileserverHits.Load()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", count)))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits counter reset to 0"))
}