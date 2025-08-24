package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/TaitA2/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
	jwtSecret      string
	apiKey         string
}

func main() {

	// Load the database
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	jwtSecret := os.Getenv("JWT_SECRET")
	apiKey := os.Getenv("POLKA_KEY")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return
	}

	const rootPath = "."
	const port = "8080"
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries:      database.New(db),
		platform:       platform,
		jwtSecret:      jwtSecret,
		apiKey:         apiKey,
	}

	serveMux := http.NewServeMux()

	handler := http.StripPrefix("/app/", http.FileServer(http.Dir(rootPath)))

	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	// OK if online
	serveMux.HandleFunc("GET /api/healthz", handlerReadiness)
	// view number of visits
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	// reset users database
	serveMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	// create user
	serveMux.HandleFunc("POST /api/users", apiCfg.handlerUsers)
	// login user
	serveMux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	// update user
	serveMux.HandleFunc("PUT /api/users", apiCfg.handlerUserUpdate)
	// create chirp
	serveMux.HandleFunc("POST /api/chirps", apiCfg.handlerChirps)
	// get all chirps
	serveMux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	// get specific chirp
	serveMux.HandleFunc("GET /api/chirps/", apiCfg.handlerGetChirp)
	// delete specific chirp
	serveMux.HandleFunc("DELETE /api/chirps/", apiCfg.handlerDeleteChirp)
	// generate new access token
	serveMux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	// revoke refresh token
	serveMux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)
	// upgrade user to chirpy red
	serveMux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerUpgrade)

	server := &http.Server{Handler: serveMux, Addr: ":" + port}
	server.ListenAndServe()
}
