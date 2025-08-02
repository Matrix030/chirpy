package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerCreateUsers(w http.ResponseWriter, r *http.Request) {
	type resquestBody struct {
		Email string `json:"email"`
	}

	var req resquestBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Email == "" {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	dbUser, err := cfg.db.CreateUser(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "could not create user", http.StatusInternalServerError)
		return
	}

	respUser := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(respUser)
}
