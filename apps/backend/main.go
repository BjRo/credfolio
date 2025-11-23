package main

import (
	"net/http"
	"os"
	"time"

	"github.com/credfolio/apps/backend/src/api"
	"github.com/credfolio/apps/backend/src/api/handlers"
	"github.com/credfolio/apps/backend/src/db"
	"github.com/credfolio/apps/backend/src/db/queries"
	"github.com/credfolio/apps/backend/src/services/extractor"
	"github.com/credfolio/apps/backend/src/services/generator"
	"github.com/credfolio/apps/backend/src/services/llm"
	"github.com/credfolio/apps/backend/src/services/profile"
	"github.com/credfolio/apps/backend/src/services/storage"
	"github.com/credfolio/apps/backend/src/utils"
	"github.com/joho/godotenv"
)

func main() {
	// Configure logger
	logger := utils.SetupLogger()

	// Load environment from .env files in development if present
	_ = godotenv.Load()
	_ = godotenv.Load(".env.development")

	port := getEnv("PORT", "8080")
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:55432/credfolio?sslmode=disable")
	openaiKey := getEnv("OPENAI_API_KEY", "")

	logger.Info("starting backend", "pid", os.Getpid(), "port", port)

	// Database
	database, err := db.Connect(dbURL)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer database.Close()

	// Repositories
	repo := queries.NewProfileQueries(database)

	// Services
	localStorage, err := storage.NewLocalStorage("uploads")
	if err != nil {
		logger.Error("failed to init storage", "error", err)
		os.Exit(1)
	}

	pdfExt := extractor.NewPDFExtractor()
	llmClient := llm.NewClient(openaiKey)
	profileExtractor := profile.NewExtractor(pdfExt, llmClient)
	cvGenerator := generator.NewCVGenerator()

	profileService := profile.NewService(repo, profileExtractor, localStorage, llmClient)

	// Handlers
	uploadHandler := handlers.NewUploadHandler(profileService)
	profileHandler := handlers.NewProfileHandler(profileService)
	cvHandler := handlers.NewCVHandler(profileService, cvGenerator)
	tailorHandler := handlers.NewTailorHandler(profileService)

	r := api.NewRouter(uploadHandler, profileHandler, cvHandler, tailorHandler)

	s := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logger.Info("backend listening", "address", ":"+port)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}

func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
