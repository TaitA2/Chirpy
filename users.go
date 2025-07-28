package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/TaitA2/Chirpy/internal/auth"
	"github.com/TaitA2/Chirpy/internal/database"
	"github.com/google/uuid"
)

type userParams struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type secureUser struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Email     string
}

func (apiCfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var params userParams
	err := decoder.Decode(&params)
	if err != nil {
		error := fmt.Sprintf("Error decoding paramters: %s", err)
		errResponse(w, error, 500)
		return
	}
	log.Printf("Login Params: %v", params)
	user, err := apiCfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		error := fmt.Sprintf("Unauthorized")
		errResponse(w, error, 401)
		return
	}
	if err = auth.CheckPasswordhash(params.Password, user.HashedPassword); err != nil {
		error := fmt.Sprintf("Unauthorized")
		errResponse(w, error, 401)
		return
	}

	data, err := json.Marshal(secureUser{user.ID, user.CreatedAt, user.UpdatedAt, user.Email})

	if err != nil {
		log.Printf("Error marshalling error response: %v", err)
	}

	w.WriteHeader(200)
	w.Write(data)
	log.Printf("Login successful!")

}

func (apiCfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var params userParams
	err := decoder.Decode(&params)
	log.Printf("Register Params: %v", params)
	if err != nil {
		error := fmt.Sprintf("Error decoding paramters: %s", err)
		errResponse(w, error, 500)
		return
	}

	// create db user with given email
	hashedPw, err := auth.HashPassword(params.Password)
	if err != nil {
		error := fmt.Sprintf("Error hashing password: %s", err)
		errResponse(w, error, 500)
		return

	}
	user, err := apiCfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPw})
	if err != nil {
		error := fmt.Sprintf("Error creating database user with email '%s': %v", params.Email, err)
		errResponse(w, error, 409)
		return
	}

	// jsonify
	data, err := json.Marshal(secureUser{user.ID, user.CreatedAt, user.UpdatedAt, user.Email})

	if err != nil {
		log.Printf("Error marshalling error response: %v", err)
	}

	w.WriteHeader(201)
	w.Write(data)

}
