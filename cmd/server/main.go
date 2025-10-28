package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"example.com/pz7-redis/internal/cache"
)

type errorResponse struct {
	Error string `json:"error"`
}

type setResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	TTL   int64  `json:"ttl_seconds"`
	Msg   string `json:"msg"`
}

type getResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ttlResponse struct {
	Key string `json:"key"`
	TTL int64  `json:"ttl_seconds"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	c := cache.New(redisAddr)
	mux := http.NewServeMux()

	mux.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		value := r.URL.Query().Get("value")
		if key == "" || value == "" {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "key and value required"})
			return
		}
		ttl := 10 * time.Second
		if err := c.Set(key, value, ttl); err != nil {
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, setResponse{
			Key:   key,
			Value: value,
			TTL:   int64(ttl / time.Second),
			Msg:   fmt.Sprintf("OK: %s=%s (TTL 10s)", key, value),
		})
	})

	mux.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if key == "" {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "key required"})
			return
		}
		val, err := c.Get(key)
		if err != nil {
			writeJSON(w, http.StatusNotFound, errorResponse{Error: err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, getResponse{Key: key, Value: val})
	})

	mux.HandleFunc("/ttl", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if key == "" {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "key required"})
			return
		}
		ttl, err := c.TTL(key)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, ttlResponse{Key: key, TTL: int64(ttl / time.Second)})
	})

	addr := "0.0.0.0:" + port
	log.Printf("Listening on %s (redis=%s)\n", addr, redisAddr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
