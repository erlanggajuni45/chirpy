package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/erlanggajuni45/chirpy/internal/auth"
	"github.com/erlanggajuni45/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		Email    string
		Password string
	}

	var req reqBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, 500, "Error when decode JSON", err)
		return
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "No token provided", err)
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "Error when getting userId", err)
		return
	}

	hashed_password, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, 500, "Error when hashing password", err)
		return
	}

	updated_user, err := cfg.database.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:          req.Email,
		HashedPassword: hashed_password,
		ID:             userID,
	})

	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, 401, "user not found", err)
			return
		}
		respondWithError(w, 500, "error when updating user", err)
		return
	}

	respondWithJSON(w, 200, User{
		ID:        updated_user.ID,
		CreatedAt: updated_user.CreatedAt,
		UpdatedAt: updated_user.UpdatedAt,
		Email:     updated_user.Email,
	})
}
