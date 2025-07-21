package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handlerValidate(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}

	type errReturn struct {
		Error string `json:"error"`
	}

	type validReturn struct {
		Valid bool `json:"valid"`
	}

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
		return
	}

	if len(params.Body) > 140 {
		log.Printf("Chirp too long")

		resp := errReturn{Error: "Chirp too long"}
		data, err := json.Marshal(resp)
		if err != nil {
			log.Printf("Error marshalling error response: %v", err)
		}

		w.WriteHeader(400)
		w.Write(data)
		return
	}
	log.Printf("Valid Chirp")
	resp := validReturn{Valid: true}
	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error marshalling error response: %v", err)
	}

	w.WriteHeader(200)
	w.Write(data)

}
