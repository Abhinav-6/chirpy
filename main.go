package main

import (
	"fmt"
	"net/http"

	"github.com/Abhinav-6/chirpy/middleware"
	"github.com/go-chi/chi/v5"
)

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("charset", "utf-8")
	fmt.Fprint(w, "OK")
}

func main() {
	var config middleware.ApiConfig
	r := chi.NewRouter()
	file := http.FileServer(http.Dir("."))
	r.Handle("/app/*", config.MiddlewareMetricsInc(http.StripPrefix("/app/", file)))
	r.Handle("/app", config.MiddlewareMetricsInc(http.StripPrefix("/app", file)))
	r.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hits: %d", config.FileServerHits)
	})

	r.Get("/reset", func(w http.ResponseWriter, r *http.Request) {
		config.FileServerHits = 0
		http.Redirect(w, r, "/metrics", http.StatusPermanentRedirect)
	})

	r.Get("/healthz", healthzHandler)
	corsMux := middleware.MiddlewareCors(r)
	http.ListenAndServe(":8080", corsMux)
}
