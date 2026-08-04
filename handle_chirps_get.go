package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	var chirpList = make([]Chirp, 0)
	chirps, err := cfg.database.GetChirps(r.Context())

	if err != nil {
		fmt.Printf("Error when getting chirps: %v\n", err)
		w.WriteHeader(500)
		w.Write([]byte("Error when getting chirps: " + err.Error()))
		return
	}

	for _, c := range chirps {
		chirpList = append(chirpList, Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(chirpList)
}
