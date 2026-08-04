package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type reqBody struct {
		Email string `json:"email"`
	}

	type errorResp struct {
		Error string `json:"error"`
	}

	var req reqBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Printf("error decode JSON: %v", err)
		w.WriteHeader(500)
		return
	}

	if req.Email == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(errorResp{
			Error: "Email can't be empty",
		})
		return
	}

	user, err := cfg.database.CreateUser(r.Context(), req.Email)
	if err != nil {
		fmt.Printf("Error when creating user: %v", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
