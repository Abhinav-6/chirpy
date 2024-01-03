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
	appRouter := chi.NewRouter()
	file := http.FileServer(http.Dir("."))
	appRouter.Handle("/*", config.MiddlewareMetricsInc(http.StripPrefix("/app/", file)))
	appRouter.Handle("/", config.MiddlewareMetricsInc(http.StripPrefix("/app", file)))

	appRouter.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hits: %d", config.FileServerHits)
	})

	appRouter.Get("/reset", func(w http.ResponseWriter, r *http.Request) {
		config.FileServerHits = 0
		http.Redirect(w, r, "/app/metrics", http.StatusPermanentRedirect)
	})

	appRouter.Get("/healthz", healthzHandler)
	r.Mount("/app/", appRouter)

	corsMux := middleware.MiddlewareCors(r)
	http.ListenAndServe(":8080", corsMux)
}
