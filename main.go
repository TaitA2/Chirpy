package main

import (
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {

	const rootPath = "."
	const port = "8080"
	var apiCfg apiConfig

	serveMux := http.NewServeMux()

	handler := http.StripPrefix("/app/", http.FileServer(http.Dir(rootPath)))

	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	serveMux.HandleFunc("GET /api/healthz", handlerReadiness)
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	serveMux.HandleFunc("POST /api/validate_chirp", handlerValidate)

	server := &http.Server{Handler: serveMux, Addr: ":" + port}
	server.ListenAndServe()
}
