package main

import (
	"log"
	"net/http"

	"bookstore/cmd"
	"bookstore/internal/config"
	"bookstore/internal/db"
)

func main() {
	cfg := config.Load()

	database := db.Connect(cfg)
	defer database.Close()

	router := cmd.SetupRoutes(database, cfg)

	log.Printf("Bookstore API starting on port %s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
