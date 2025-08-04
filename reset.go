package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.env != "dev" {
		http.Error(w, "forbidden in non-dev environment", http.StatusForbidden)
		return
	}

	cfg.fileserverHits.Store(0)

	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		http.Error(w, "failed to reset users", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0 and users deleted"))
}
