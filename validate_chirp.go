package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handleChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type requestBody struct {
		Body string `json:"body"`
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
	respondWithJson(w, http.StatusOK, cleaned_body)
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

func profanityChecker(s string) map[string]string {
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
	return map[string]string{"cleaned_body": cleaned}
}
