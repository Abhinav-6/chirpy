package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"io"
	// "log"
	"net/http"

	"github.com/Abhinav-6/chirpy/database"
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
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<html>

			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
			
			</html>`, config.FileServerHits)
	})

	r.Mount("/app/", appRouter)
	r.Mount("/admin", adminRouter)

	r.Post("/api/validate_chirp", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		type param struct {
			Body string `json:"body"`
		}
		type errorResponse struct {
			Message string `json:"Error"`
		}
		type validResponse struct {
			Message string `json:"Valid"`
		}
		dat, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "couldn't read request")
			return
		}
		var params param
		err = json.Unmarshal(dat, &params)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "error parsing request")
			return
		}

		if len(params.Body) > 140 {
			w.WriteHeader(400)
			fmt.Fprintf(w, errorResponse{"Chirp is too long"}.Message)
			return
		}
		out, err := json.Marshal(validResponse{cleanChirp(params.Body)})
		if err != nil {
			panic(err)
		}
		w.WriteHeader(200)
		fmt.Fprintf(w, string(out))

	})

	corsMux := middleware.MiddlewareCors(r)
	http.ListenAndServe(":3000", corsMux)
}

func cleanChirp(s string) string {
	ss := strings.Split(s, " ")
	for i, str := range ss {
		word := strings.ToLower(str)
		switch word {
		case "fornax", "sharbert", "kerfuffle":
			ss[i] = "****"
		}
	}
	return strings.Join(ss, " ")
}
