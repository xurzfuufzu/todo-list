package main

import (
	"log"
	"net/http"
	"todoList/internal/config"
	"todoList/internal/router"
)

func main() {
	cfg := config.MustLoad()

	router := router.NewRouter()
	log.Printf("Server is running at %s", cfg.HttpServer.Port)
	if err := http.ListenAndServe(cfg.HttpServer.Port, router); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
