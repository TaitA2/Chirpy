package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

func main() {
	serveMux := http.NewServeMux()
	server := &http.Server{Handler: serveMux, Addr: ":8080"}

	var apiCfg apiConfig
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	serveMux.HandleFunc("/healthz/", handlerReadiness)
	serveMux.HandleFunc("/metrics/", func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
		writer.WriteHeader(200)
		writer.Write(fmt.Appendf(nil, "Hits: %v", apiCfg.fileserverHits.Load()))
	})

	serveMux.HandleFunc("/reset/", func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
		writer.WriteHeader(200)
		writer.Write([]byte("OK"))
		apiCfg.fileserverHits.Swap(0)
	})

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

func handlerReadiness(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write([]byte("OK"))
}
