package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/erlanggajuni45/chirpy/internal/auth"
)

type RefreshTokenRes struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	refresh_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		fmt.Println("Error get refresh token: ", err.Error())
		respondWithError(w, 500, "Error get refresh token", err)
		return
	}

	user_id, err := cfg.database.GetUserFromRefreshToken(r.Context(), refresh_token)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(401)
			w.Write([]byte("Refresh token is invalid"))
			return
		}
		fmt.Println(err.Error())
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	fmt.Println("INIKAH", user_id)

	token, err := auth.MakeJWT(user_id, cfg.jwtSecret, time.Duration(1)*time.Hour)

	if err != nil {
		fmt.Println("Error when get jwt token: ", err.Error())
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	respondWithJSON(w, 200, RefreshTokenRes{
		Token: token,
	})
}
