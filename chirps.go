package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/TaitA2/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(strings.TrimPrefix(r.URL.Path, "/api/chirps/"))
	if err != nil {
		error := fmt.Sprintf("Error parsing url: %v", err)
		errResponse(w, error)
		return
	}
	log.Printf("chirp id: %v", chirpID)
	chirp, err := apiCfg.dbQueries.GetChirp(r.Context(), chirpID)
	if err != nil {
		error := fmt.Sprintf("Error getting chirp: %v", err)
		errResponse(w, error)
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
	chirp, err := validateChirp(w, r)
	if err != nil {
		log.Printf("Error validating chrip: %v", err)
		return
	}

	dbChrip, err := apiCfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   chirp.CleanedBody,
		UserID: chirp.UserID,
	})
	if err != nil {
		error := fmt.Sprintf("Error creating database entry for chirp: %v", err)
		errResponse(w, error)
		return
	}

	resp, err := json.Marshal(dbChrip)
	w.WriteHeader(201)
	w.Write(resp)

}

func errResponse(w http.ResponseWriter, s string) {
	log.Printf("%s", s)
	resp := errReturn{Error: s}
	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error marshalling error response: %v", err)
	}

	w.WriteHeader(500)
	w.Write(data)
}
