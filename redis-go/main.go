package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
	"redis-go/storage"
	"time"
)

func main() {
	cfg := storage.Config{
		Addr:        "localhost:6379",
		Password:    "test1234",
		User:        "testuser",
		DB:          0,
		MaxRetry:    5,
		DialTimeout: 10 * time.Second,
		Timeout:     5 * time.Second,
	}

	db, err := storage.NewCLient(context.Background(), cfg)
	if err != nil {
		panic(err)
	}

	router := chi.NewRouter()
	router.Route("/card", storage.NewCardHandler(context.Background(), db))
	srv := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
