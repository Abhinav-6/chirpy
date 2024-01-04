package main

import (
	// "encoding/json"
	"fmt"
	"strings"

	// "io"
	// "log"
	"github.com/go-chi/chi/middleware"
	"net/http"

	"github.com/Abhinav-6/chirpy/assets/util"
	"github.com/Abhinav-6/chirpy/database"
	// "github.com/Abhinav-6/chirpy/middleware"
	"github.com/go-chi/chi/v5"
)

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("charset", "utf-8")
	fmt.Fprint(w, "OK")
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RedirectSlashes)
	db, _ := database.NewDb()

	r.Get("/api/chirps", func(w http.ResponseWriter, r *http.Request) {
		data, err := db.GetChirps()
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		util.RespondWithJSON(w, http.StatusOK, data)
	})

	r.Get("/api/chirps/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		data, err := db.GetChirpsById(id)
		if err != nil {
			util.RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		util.RespondWithJSON(w, http.StatusOK, data)
	})

	// r.Post("/api/validate_chirp", func(w http.ResponseWriter, r *http.Request) {
	// 	defer r.Body.Close()

	// 	dat, err := io.ReadAll(r.Body)
	// 	if err != nil {
	// 		w.WriteHeader(500)
	// 		fmt.Fprintf(w, "couldn't read request")
	// 		return
	// 	}
	// 	var params param
	// 	err = json.Unmarshal(dat, &params)
	// 	if err != nil {
	// 		w.WriteHeader(500)
	// 		fmt.Fprintf(w, "error parsing request")
	// 		return
	// 	}

	// 	if len(params.Body) > 140 {
	// 		w.WriteHeader(400)
	// 		fmt.Fprintf(w, errorResponse{"Chirp is too long"}.Message)
	// 		return
	// 	}
	// 	out, err := json.Marshal(validResponse{cleanChirp(params.Body)})
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	w.WriteHeader(200)
	// 	fmt.Fprintf(w, string(out))

	// })

	corsMux := util.MiddlewareCors(r)
	http.ListenAndServe("localhost:3000", corsMux)
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
