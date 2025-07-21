package main

import "net/http"

func (apiCfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	isValid := handlerValidate(w, r)

}
