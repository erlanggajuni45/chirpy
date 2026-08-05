package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/erlanggajuni45/chirpy/internal/auth"
)

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

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
