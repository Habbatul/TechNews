package main

import (
	"TechNews/config"
	"TechNews/memory"
	"log"
	"net/http"
)

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
	log.Fatal(http.ListenAndServe(":8080", mux))
}
