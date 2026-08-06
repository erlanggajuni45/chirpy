package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/erlanggajuni45/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author_id")

	var chirps = make([]database.Chirp, 0)
	var err error

	if author != "" {
		userId, err := uuid.Parse(author)
		if err != nil {
			fmt.Printf("Error when parsing to UUID: %v\n", err)
			w.WriteHeader(500)
			return
		}

		chirps, err = cfg.database.GetChirpsByUser(r.Context(), userId)
		if err != nil {
			fmt.Printf("Error when getting chirps by user id: %v\n", err)
			w.WriteHeader(500)
			return
		}
	} else {
		chirps, err = cfg.database.GetChirps(r.Context())

		if err != nil {
			fmt.Printf("Error when getting chirps: %v\n", err)
			w.WriteHeader(500)
			w.Write([]byte("Error when getting chirps: " + err.Error()))
			return
		}
	}

	sortQuery := r.URL.Query().Get("sort")
	if sortQuery == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	}

	var chirpList = make([]Chirp, 0)
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

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		fmt.Printf("Error when parsing to UUID: %v\n", err)
		w.WriteHeader(500)
		return
	}

	chirp, err := cfg.database.GetChirp(r.Context(), id)
	if err != nil {
		fmt.Printf("Error when getting chirp: %v\n", err)
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			w.Write([]byte("chirp not found"))
			return
		}

		w.WriteHeader(500)
		w.Write([]byte("Error when getting chirp: " + err.Error()))
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
