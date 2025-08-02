package main

import (
	"database/sql"
	"github.com/Matrix030/chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)
import _ "github.com/lib/pq"

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	queries := database.New(db)

	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		db:             queries,
	}

	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", handleChirp)
	mux.HandleFunc("POST /api/users", apiCfg.handleCreateUser)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
