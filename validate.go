package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func validateChirp(w http.ResponseWriter, r *http.Request) (validReturn, error) {
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	log.Printf("Params: %v", params)
	if err != nil {
		log.Printf("Error decoding paramters: %s", err)

		resp := errReturn{Error: "Something went wrong"}
		data, err := json.Marshal(resp)
		if err != nil {
			log.Printf("Error marshalling error response: %v", err)
		}

		w.WriteHeader(500)
		w.Write(data)
		return validReturn{}, err
	}

	// Ensure Chirp does not exceed 140 character limit
	if len(params.Body) > 140 {
		log.Printf("Chirp too long")

		resp := errReturn{Error: "Chirp too long"}
		data, err := json.Marshal(resp)
		if err != nil {
			log.Printf("Error marshalling error response: %v", err)
		}

		w.WriteHeader(400)
		w.Write(data)
		return validReturn{}, err
	}

	// Profanity filter
	cleanedBody := cleanBody(params.Body)
	log.Printf("Censored: %s", params.Body)

	// Chirp is valid
	log.Printf("Valid Chirp")
	resp := validReturn{
		CleanedBody: cleanedBody,
		UserID:      params.UserID,
	}

	return resp, nil

}

func cleanBody(body string) string {
	const censor = "****"
	profanes := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(body, " ")
	for i := range words {
		for j := range profanes {
			if strings.ToLower(words[i]) == profanes[j] {
				words[i] = censor
			}
		}
	}
	return strings.Join(words, " ")
}

type parameters struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type errReturn struct {
	Error string `json:"error"`
}

type validReturn struct {
	CleanedBody string    `json:"cleaned_body"`
	UserID      uuid.UUID `json:"user_id"`
}
