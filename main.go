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
	adminRouter := chi.NewRouter()
	file := http.FileServer(http.Dir("."))
	appRouter.Handle("/*", config.MiddlewareMetricsInc(http.StripPrefix("/app/", file)))
	appRouter.Handle("/", config.MiddlewareMetricsInc(http.StripPrefix("/app", file)))

	adminRouter.Get("/reset", func(w http.ResponseWriter, r *http.Request) {
		config.FileServerHits = 0
		http.Redirect(w, r, "/admin/metrics", http.StatusPermanentRedirect)
	})

	appRouter.Get("/healthz", healthzHandler)

	adminRouter.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type","text/html")
		fmt.Fprintf(w,`<html>

			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
			
			</html>`, config.FileServerHits)
	})

	r.Mount("/app/", appRouter)
	r.Mount("/admin", adminRouter)

	corsMux := middleware.MiddlewareCors(r)
	http.ListenAndServe(":8080", corsMux)
}
