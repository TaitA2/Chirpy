package main

import (
	"context"
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
	ID           uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Email        string
	Token        string
	RefreshToken string
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

	AccessJWT, err := auth.MakeJWT(user.ID, apiCfg.jwtSecret)
	if err != nil {
		error := fmt.Sprintf("Error making JWT.: %v", err)
		errResponse(w, error, 401)
		return
	}

	RefreshToken, err := auth.MakeRefreshToken()
	log.Printf("Refresh Token: %v", RefreshToken)
	if err != nil {
		error := fmt.Sprintf("Error making refresh token: %v", err)
		errResponse(w, error, 401)
		return
	}
	apiCfg.dbQueries.CreateToken(context.Background(), database.CreateTokenParams{
		Token:     RefreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(time.Duration(time.Hour * 24 * 60)),
	})

	data, err := json.Marshal(secureUser{user.ID, user.CreatedAt, user.UpdatedAt, user.Email, AccessJWT, RefreshToken})

	if err != nil {
		log.Printf("Error marshalling error response: %v", err)
	}

	w.WriteHeader(200)
	w.Write(data)
	log.Printf("Login Successful!")

}

func (apiCfg *apiConfig) handlerUserUpdate(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		errResponse(w, fmt.Sprintf("Error getting bearer token: %v", err), 401)
		return
	}

	userID, err := auth.ValidateJWT(refreshToken, apiCfg.jwtSecret)
	if err != nil {
		errResponse(w, fmt.Sprintf("%v", err), 401)
		return
	}
	// get new user params
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var params userParams
	err = decoder.Decode(&params)
	if err != nil {
		error := fmt.Sprintf("Error decoding paramters: %s", err)
		errResponse(w, error, 500)
		return
	}

	// hash pw
	hashedPw, err := auth.HashPassword(params.Password)
	if err != nil {
		error := fmt.Sprintf("Error hashing password: %s", err)
		errResponse(w, error, 500)
		return

	}

	// update user
	err = apiCfg.dbQueries.UpdateUser(context.Background(), database.UpdateUserParams{ID: userID, Email: params.Email, HashedPassword: hashedPw})
	if err != nil {
		errResponse(w, fmt.Sprintf("Error updating user: %v", err), 500)
		return
	}

	// respond
	user, err := apiCfg.dbQueries.GetUserByEmail(context.Background(), params.Email)
	if err != nil {
		errResponse(w, fmt.Sprintf("Error getting user by updated email: %v", err), 400)
		return
	}
	accessJWT, err := auth.MakeJWT(user.ID, apiCfg.jwtSecret)
	data, err := json.Marshal(secureUser{user.ID, user.CreatedAt, user.UpdatedAt, user.Email, accessJWT, refreshToken})

	if err != nil {
		log.Printf("Error marshalling error response: %v", err)
		return
	}

	w.WriteHeader(200)
	w.Write(data)
	log.Printf("Login Successful!")

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
	tokenString, err := auth.MakeJWT(user.ID, apiCfg.jwtSecret)
	if err != nil {
		error := fmt.Sprintf("Error making JWT.: %v", err)
		errResponse(w, error, 401)
		return
	}
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		error := fmt.Sprintf("Error making refresh token: %v", err)
		errResponse(w, error, 401)
		return
	}
	apiCfg.dbQueries.CreateToken(context.Background(), database.CreateTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(time.Duration(time.Hour * 24 * 60)),
	})
	data, err := json.Marshal(secureUser{user.ID, user.CreatedAt, user.UpdatedAt, user.Email, tokenString, refreshToken})

	if err != nil {
		log.Printf("Error marshalling error response: %v", err)
	}

	w.WriteHeader(201)
	w.Write(data)

}
