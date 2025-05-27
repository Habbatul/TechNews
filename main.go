package main

import (
	"TechNews/config"
	"TechNews/memory"
	_ "embed"
	"log"
	"net/http"
)

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func main() {
	config.InitRedis()

	//routes
	mux := http.NewServeMux()
	mux.HandleFunc("/api/feed", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			memory.GetFeedsFromRedis(w)
		} else if r.Method == "POST" {
			memory.SaveFeedsData(w)
		}
	})

	mux.HandleFunc("/api/resume", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			memory.GetResumeFromRedis(w)
		} else if r.Method == "POST" {
			memory.SaveResumeData(w)
		}
	})

	//miiddleware (buat buka CORS)
	handler := withCORS(mux)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
