package main

import (
	"net/http"

	"github.com/erlanggajuni45/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpId, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, 500, "error when getting chirp id", err)
		return
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "error when getting access token", err)
		return
	}

	userId, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "error when getting user id", err)
		return
	}

	chirp, err := cfg.database.GetChirp(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, 500, "error when getting chirp", err)
		return
	}

	if chirp.UserID != userId {
		respondWithError(w, 403, "can't delete other's chirp", err)
		return
	}

	err = cfg.database.DeleteChirp(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, 500, "error when deleting chirp", err)
		return
	}

	w.WriteHeader(204)
}
