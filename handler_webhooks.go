package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/erlanggajuni45/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerWebhook(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 401, "error when getting api key", err)
		return
	}

	if apiKey != cfg.apiKey {
		respondWithError(w, 401, "api key didn't match!", err)
		return
	}

	type dataReq struct {
		UserId string `json:"user_id"`
	}

	type parameters struct {
		Event string  `json:"event"`
		Data  dataReq `json:"data"`
	}

	var req parameters
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, 500, "error when decode payload", err)
		return
	}

	userId, err := uuid.Parse(req.Data.UserId)
	if err != nil {
		respondWithError(w, 500, "error when parsing user id", err)
		return
	}

	switch req.Event {
	case "user.upgraded":
		user, err := cfg.database.UpgradeUser(r.Context(), userId)
		if err != nil {
			if err == sql.ErrNoRows {
				respondWithError(w, 404, "user not found", err)
				return
			}
			respondWithError(w, 500, "error when upgrading user", err)
			return
		}
		w.WriteHeader(204)
		json.NewEncoder(w).Encode(User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed.Bool,
		})
	default:
		w.WriteHeader(204)
	}
}
