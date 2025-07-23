package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/TaitA2/Chirpy/internal/database"
)

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
		log.Printf("Want: %v, Have: %v", chirp.UserID, dbChrip.UserID)
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
