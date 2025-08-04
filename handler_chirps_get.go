package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}

	chirps := []Chirps{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, chirpStruct{
			ID:       dbChirp.ID,
			CreateAt: dbChirp.CreatedAt,
			UpdateAt: dbChirp.UpdatedAt,
			UserID:   dbChirp.UserID,
			Body:     dbChirp.Body,
		})
	}

	respondWithJson(w, http.StatusOK, chirps)
}
