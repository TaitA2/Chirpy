package main

import (
	"log"
	"net/http"
)

func (apiCfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	if apiCfg.platform == "dev" {
		w.WriteHeader(200)
		w.Write([]byte("OK"))
		apiCfg.fileserverHits.Swap(0)
		if err := apiCfg.dbQueries.ResetUsers(r.Context()); err != nil {
			log.Printf("Error resetting users database: %v", err)
		}
		return
	}
	w.WriteHeader(403)
	w.Write([]byte("Forbidden"))
}
