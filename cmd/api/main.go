package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/alternative/backend/internal/platform/auth"
	"github.com/alternative/backend/internal/platform/database"
	"github.com/alternative/backend/internal/platform/seed"
	"github.com/alternative/backend/internal/platform/server"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
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

	// Auto-seed
	seedAdmin()
	seed.SeedMockListings()

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

func seedAdmin() {
	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		adminEmail = "admin@alternative.lat"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	coll := database.DB.Collection("users")
	result, err := coll.UpdateOne(ctx,
		bson.M{"email": adminEmail, "role": bson.M{"$ne": "admin"}},
		bson.M{"$set": bson.M{"role": "admin"}},
	)
	if err != nil {
		log.Printf("Admin seed check failed: %v", err)
		return
	}
	if result.ModifiedCount > 0 {
		log.Printf("Seeded admin role for %s", adminEmail)
	}
}
