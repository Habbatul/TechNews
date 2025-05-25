package memory

import (
	"TechNews/config"
	"TechNews/service"
	"encoding/json"
	"net/http"
)

func SaveFeedsData(w http.ResponseWriter) {
	jsonData, err := json.Marshal(service.GetNews())
	if err != nil {
		http.Error(w, "Something wrong", http.StatusInternalServerError)
		return
	}

	err = config.RedisClient.Set(config.Ctx, "feeds", jsonData, 0).Err()
	if err != nil {
		http.Error(w, "Failed to save to Redis", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "data stored in Redis"}`))
}

func SaveResumeData(w http.ResponseWriter) {
	jsonData, err := json.Marshal(service.GetResume())
	if err != nil {
		http.Error(w, "Something wrong", http.StatusInternalServerError)
		return
	}

	err = config.RedisClient.Set(config.Ctx, "resume", jsonData, 0).Err()
	if err != nil {
		http.Error(w, "Failed to save to Redis", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "news resume stored in Redis"}`))
}

func GetFeedsFromRedis(w http.ResponseWriter) {
	data, err := config.RedisClient.Get(config.Ctx, "feeds").Result()
	if err != nil {
		http.Error(w, "Data not found in Redis", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(data))
}

func GetResumeFromRedis(w http.ResponseWriter) {
	data, err := config.RedisClient.Get(config.Ctx, "resume").Result()
	if err != nil {
		http.Error(w, "Resume not found in Redis", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(data))
}
