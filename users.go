package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (apiCfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params struct {
		Email string `json:"email"`
	}
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
		return
	}

	// create db user with given email
	user, err := apiCfg.dbQueries.CreateUser(r.Context(), params.Email)
	if err != nil {
		log.Printf("Error creating database user with email '%s': %v", params.Email, err)
		resp := errReturn{Error: "Error creating user, email already in use."}
		data, err := json.Marshal(resp)
		if err != nil {
			log.Printf("Error marshalling error response: %v", err)
		}

		w.WriteHeader(500)
		w.Write(data)
		return
	}

	// jsonify
	data, err := json.Marshal(user)
	if err != nil {
		log.Printf("Error marshalling error response: %v", err)
	}

	w.WriteHeader(201)
	w.Write(data)

}
