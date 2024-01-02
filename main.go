package main

import (
	"fmt"
	"net/http"

	"github.com/Abhinav-6/chirpy/middleware"
)

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("charset", "utf-8")
	fmt.Fprint(w, "OK")
}

func main() {
	var config middleware.ApiConfig
	mux := http.NewServeMux()
	file := http.FileServer(http.Dir("."))
	mux.Handle("/app/", config.MiddlewareMetricsInc(http.StripPrefix("/app/", file)))

	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hits: %d", config.FileServerHits)
	})

	mux.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		config.FileServerHits=0;
		http.Redirect(w,r,"/metrics",http.StatusPermanentRedirect)
	})

	mux.HandleFunc("/healthz", healthzHandler)
	corsMux := middleware.MiddlewareCors(mux)
	http.ListenAndServe(":8080", corsMux)
}
