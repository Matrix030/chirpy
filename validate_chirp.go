package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

type chirpStruct struct {
	ID       uuid.UUID `json:"id"`
	CreateAt time.Time `json:"created_at"`
	UpdateAt time.Time `json:"updated_at"`
	Body     string    `json:"body"`
	UserID   uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handleChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type requestBody struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
	}

	chirpReq := &requestBody{}
	err := json.NewDecoder(r.Body).Decode(&chirpReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if len(chirpReq.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	cleaned_body := profanityChecker(chirpReq.Body)
	final_body := &chirpStruct{
		ID:       uuid.New(),
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
		Body:     cleaned_body,
		UserID:   chirpReq.UserID,
	}
	respondWithJson(w, http.StatusCreated, final_body)
}

func profanityChecker(s string) string {
	profanityWords := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}

	words := strings.Split(s, " ")
	for i, word := range words {
		lowerWord := strings.ToLower(word)
		if profanityWords[lowerWord] {
			words[i] = "****"
		}
	}

	cleaned := strings.Join(words, " ")
	return cleaned
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	return respondWithJson(w, code, map[string]string{"error": msg})
}
