package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/erlanggajuni45/chirpy/internal/auth"
)

type RefreshTokenRes struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	refresh_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		fmt.Println("Error get refresh token: ", err.Error())
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
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

	token, err := auth.MakeJWT(user_id, cfg.jwtSecret, time.Duration(1)*time.Hour)

	if err != nil {
		fmt.Println("Error when get jwt token: ", err.Error())
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(RefreshTokenRes{
		Token: token,
	})
}
