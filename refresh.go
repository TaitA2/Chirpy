package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/TaitA2/Chirpy/internal/auth"
)

func (apiCfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		errResponse(w, fmt.Sprintf("Error getting bearer: %v", err), 401)
		return
	}
	dbToken, err := apiCfg.dbQueries.GetToken(context.Background(), token)
	if err != nil {
		errResponse(w, "Error getting token: refresh token not found", 401)
		return
	}
	if dbToken.ExpiresAt.Before(time.Now()) {
		errResponse(w, "Error getting token: refresh token is expired", 401)
		return
	}
	if dbToken.RevokedAt.Valid {
		errResponse(w, "Error getting token: refresh token is revoked", 401)
		return

	}

	accessToken, err := auth.MakeJWT(dbToken.UserID, apiCfg.jwtSecret)
	data, err := json.Marshal(struct {
		Token string `json:"token"`
	}{accessToken})
	w.WriteHeader(200)
	w.Write(data)
	log.Printf("Refresh Successful!\nToken: %v", accessToken)
	return
}
