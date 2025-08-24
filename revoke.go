package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/TaitA2/Chirpy/internal/auth"
)

func (apiCfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		errResponse(w, fmt.Sprintf("Errror getting refresh token: %v", err), 401)
		return
	}
	if err = apiCfg.dbQueries.RevokeToken(context.Background(), token); err != nil {
		errResponse(w, fmt.Sprintf("Error revoking token: %v", err), 401)
		return
	}
	log.Printf("Token revoked successfully!")
	w.WriteHeader(204)
	return
}
