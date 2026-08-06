package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/erlanggajuni45/chirpy/internal/auth"
	"github.com/erlanggajuni45/chirpy/internal/database"
	"github.com/google/uuid"
)

type LoginResp struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req params
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println("error when decode JSON:", err.Error())
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	user, err := cfg.database.GetUser(r.Context(), req.Email)
	if err != nil {
		fmt.Println("error when getting user: ", err.Error())
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	is_match, err := auth.CheckPasswordHash(req.Password, user.HashedPassword)
	if err != nil {
		fmt.Println("error when checking password: ", err.Error())
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	if !is_match {
		w.WriteHeader(401)
		w.Write([]byte("Incorrect email or password"))
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)

	if err != nil {
		fmt.Println("Error when get jwt token: ", err.Error())
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	refresh_token := auth.MakeRefreshToken()
	err = cfg.database.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		UserID: user.ID,
		Token:  refresh_token,
	})

	if err != nil {
		fmt.Println("Error when storing refresh token")
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(LoginResp{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refresh_token,
		IsChirpyRed:  user.IsChirpyRed.Bool,
	})
}
