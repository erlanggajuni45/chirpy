package main

import (
	"fmt"
	"net/http"

	"github.com/erlanggajuni45/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshTokenRevoke(w http.ResponseWriter, r *http.Request) {
	refresh_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		fmt.Println("Error get refresh token: ", err.Error())
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	err = cfg.database.RevokeRefreshToken(r.Context(), refresh_token)
	if err != nil {
		fmt.Println("Error when revoking refresh token:", err.Error())
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(204)
}
