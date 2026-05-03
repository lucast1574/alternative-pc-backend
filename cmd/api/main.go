package main

import (
	"log"
	"os"

	"github.com/alternative/backend/internal/platform/database"
	"github.com/alternative/backend/internal/platform/server"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	_, err := database.Connect(mongoURI)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	srv := server.NewServer()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := srv.Run(port); err != nil {
		log.Fatalf("Could not run server: %v", err)
	}
}
