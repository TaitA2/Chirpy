package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/TaitA2/Chirpy/internal/auth"
	"github.com/TaitA2/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(strings.TrimPrefix(r.URL.Path, "/api/chirps/"))
	if err != nil {
		error := fmt.Sprintf("Error parsing url: %v", err)
		errResponse(w, error, 500)
		return
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		errResponse(w, fmt.Sprintf("%v", err), 401)
		return
	}
	userID, err := auth.ValidateJWT(token, apiCfg.jwtSecret)
	if err != nil {
		errResponse(w, fmt.Sprintf("%v", err), 401)
		return
	}
	err = apiCfg.dbQueries.DeleteChirp(context.Background(), database.DeleteChirpParams{UserID: userID, ID: chirpID})
	if err != nil {
		errResponse(w, fmt.Sprintf("Error deleting chirp: %v", err), 500)
		return
	}
	w.WriteHeader(204)
	log.Printf("Chirp deleted successfully!")
}

func (apiCfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(strings.TrimPrefix(r.URL.Path, "/api/chirps/"))
	if err != nil {
		error := fmt.Sprintf("Error parsing url: %v", err)
		errResponse(w, error, 500)
		return
	}
	log.Printf("chirp id: %v", chirpID)
	chirp, err := apiCfg.dbQueries.GetChirp(r.Context(), chirpID)
	if err != nil {
		error := fmt.Sprintf("Error getting chirp: %v", err)
		errResponse(w, error, 404)
		return
	}
	data, err := json.Marshal(chirp)
	if err != nil {
		log.Printf("Error marshalling error response: %v", err)
	}

	w.WriteHeader(200)
	w.Write(data)
	return

}
func (apiCfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {

	chirps, err := apiCfg.dbQueries.GetChirps(r.Context())
	if err != nil {
		log.Printf("Error getting chirps: %v", err)
		resp := errReturn{Error: "Error getting chirps"}
		data, err := json.Marshal(resp)
		if err != nil {
			log.Printf("Error marshalling error response: %v", err)
		}

		w.WriteHeader(500)
		w.Write(data)
		return
	}
	resp := []database.Chirp{}
	for i := range chirps {
		resp = append(resp, chirps[i])
	}
	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error marshalling error response: %v", err)
	}

	w.WriteHeader(200)
	w.Write(data)
	return

}

func (apiCfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {

	// auth
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		errResponse(w, fmt.Sprintf("%v", err), 401)
		return
	}
	userID, err := auth.ValidateJWT(tokenString, apiCfg.jwtSecret)
	if err != nil {
		errResponse(w, fmt.Sprintf("%v", err), 401)
		return
	}

	chirp, err := validateChirp(w, r)
	if err != nil {
		log.Printf("Error validating chrip: %v", err)
		return
	}

	dbChrip, err := apiCfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   chirp.CleanedBody,
		UserID: userID,
	})
	if err != nil {
		error := fmt.Sprintf("Error creating database entry for chirp: %v", err)
		errResponse(w, error, 500)
		return
	}

	resp, err := json.Marshal(dbChrip)
	w.WriteHeader(201)
	w.Write(resp)

}

func errResponse(w http.ResponseWriter, s string, code int) {
	log.Printf("%s", s)
	resp := errReturn{Error: s}
	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error marshalling error response: %v", err)
	}

	w.WriteHeader(code)
	w.Write(data)
}
