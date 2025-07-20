package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func main() {

	var apiCfg apiConfig

	serveMux := http.NewServeMux()

	handler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))

	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	serveMux.HandleFunc("GET /api/healthz", handlerReadiness)
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	serveMux.HandleFunc("POST /api/validate_chirp", handlerValidate)

	server := &http.Server{Handler: serveMux, Addr: ":8080"}
	server.ListenAndServe()
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func handlerReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (apiCfg *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(200)
	w.Write(fmt.Appendf(nil, `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, apiCfg.fileserverHits.Load()))
}

func (apiCfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
	apiCfg.fileserverHits.Swap(0)
}

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
