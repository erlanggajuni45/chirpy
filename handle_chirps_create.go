package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/erlanggajuni45/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	type errorResp struct {
		Error string `json:"error"`
	}

	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Printf("Error decode json: %v\n", err)
		w.WriteHeader(500)
		return
	}

	cleaned_body, err := validateChirp(params.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	new_chirp, err := cfg.database.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleaned_body,
		UserID: params.UserId,
	})

	if err != nil {
		fmt.Printf("error when inserting chirp: %v\n", err)
		w.WriteHeader(500)
		w.Write([]byte("Error inserting chirp: " + err.Error()))
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(Chirp{
		ID:        new_chirp.ID,
		CreatedAt: new_chirp.CreatedAt,
		UpdatedAt: new_chirp.UpdatedAt,
		Body:      new_chirp.Body,
		UserID:    new_chirp.UserID,
	})
}

func validateChirp(body string) (string, error) {
	if len(body) > 140 {
		return "", errors.New("Chirp is too long")
	}

	return getCleanedBody(body), nil
}

func getCleanedBody(body string) string {
	list_word := []string{}
	filter_word := []string{"kerfuffle", "sharbert", "fornax"}
	for _, word := range strings.Split(body, " ") {
		selected_word := word
		if slices.Contains(filter_word, strings.ToLower(word)) {
			selected_word = "****"
		}
		list_word = append(list_word, selected_word)
	}

	return strings.Join(list_word, " ")
}
