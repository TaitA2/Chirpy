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
}

func main() {

	// Load the database
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
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
		platform:       platform}

	serveMux := http.NewServeMux()

	handler := http.StripPrefix("/app/", http.FileServer(http.Dir(rootPath)))

	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	serveMux.HandleFunc("GET /api/healthz", handlerReadiness)
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	serveMux.HandleFunc("POST /api/users", apiCfg.handlerUsers)
	serveMux.HandleFunc("POST /api/chirps", apiCfg.handlerChirps)
	serveMux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)

	server := &http.Server{Handler: serveMux, Addr: ":" + port}
	server.ListenAndServe()
}
