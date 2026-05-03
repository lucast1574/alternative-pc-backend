package main

import (
	"log"
	"os"

	"github.com/alternative/backend/internal/platform/auth"
	"github.com/alternative/backend/internal/platform/database"
	"github.com/alternative/backend/internal/platform/server"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	// Init JWT
	auth.Init()

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "alternative-pc"
	}

	_, err := database.Connect(mongoURI, dbName)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	srv := server.NewServer()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := srv.Run(port); err != nil {
		log.Fatalf("Could not run server: %v", err)
	}
}
