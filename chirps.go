package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/TaitA2/Chirpy/internal/database"
)

func (apiCfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := apiCfg.dbQueries.GetChirps(r.Context())
	if err != nil {
		log.Printf("Error getting chirps: %v", err)
		resp := errReturn{Error: "Something getting chirps"}
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

	w.WriteHeader(500)
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
		log.Printf("Error creating database entry for chirp: %v", err)
		resp := errReturn{Error: "Something went wrong"}
		data, err := json.Marshal(resp)
		if err != nil {
			log.Printf("Error marshalling error response: %v", err)
		}

		w.WriteHeader(500)
		w.Write(data)
		return
	}

	resp, err := json.Marshal(dbChrip)
	w.WriteHeader(201)
	w.Write(resp)

}
