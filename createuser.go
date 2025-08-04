package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUsers(w http.ResponseWriter, r *http.Request) {
	type resquestBody struct {
		Email string `json:"email"`
	}
	defer r.Body.Close()

	var req resquestBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Email == "" {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	// create the user in the database using go generated method
	dbUser, err := cfg.db.CreateUser(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "could not create user", http.StatusInternalServerError)
		return
	}
	// this is the response the user will get once everything is executed correctly
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
