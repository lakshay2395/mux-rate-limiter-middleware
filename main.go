package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
	"github.com/lakshay2395/mux-rate-limiter-middleware/middleware"
)

func main() {
	client, _ := getRedisClient()
	limiter := middleware.LeakyBucket(client, 10, 60*time.Second, func(r *http.Request) string {
		return "U2"
	})
	r := mux.NewRouter()
	r.Use(limiter)
	r.HandleFunc("/ok", OkHandler)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", "8090"), r))
}

func getRedisClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func OkHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}
