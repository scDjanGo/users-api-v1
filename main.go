package main

import (
	"log"
	"strings"
	"time"

	"github.com/joho/godotenv"

	_ "users-api-v1/src/docs"
	"users-api-v1/src/features"
	"users-api-v1/src/handlers"
	"users-api-v1/src/media"
	"users-api-v1/src/repository"
	"users-api-v1/src/router"
	"users-api-v1/src/seed"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Printf(".env не найден, использую переменные окружения: %v", err)
	}

	mediaDir := features.GetEnv("MEDIA_DIR", "./data/media")
	baseURL := strings.TrimRight(features.GetEnv("BASE_URL", "http://localhost:8080"), "/")
	port := features.GetEnv("PORT", "8080")

	storage, err := media.NewStorage(mediaDir)

	if err != nil {
		log.Fatalf("Media storage: %v", err)
	}

	// Connect DB
	userRepo, err := repository.NewUserRepo()
	if err != nil {
		log.Fatalf("open db: %v", err)
	}

	if err := seed.Seed(userRepo); err != nil {
		log.Fatalf("seed: %v", err)
	}
	seed.StartPeriodic(userRepo, 3*24*time.Hour)

	userHandler := handlers.NewUserHandler(userRepo, storage, baseURL)

	r := router.RouterInit(userHandler, mediaDir, storage.MaxSize())
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server: %v", err)
	}
}
