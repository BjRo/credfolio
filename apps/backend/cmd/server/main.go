package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/credfolio/apps/backend/internal/handler/middleware"
	"github.com/credfolio/apps/backend/internal/repository"
	"github.com/credfolio/apps/backend/pkg/config"
	"github.com/credfolio/apps/backend/pkg/logger"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logLevel := logger.LevelInfo
	if cfg.IsDevelopment() {
		logLevel = logger.LevelDebug
	}
	l := logger.New(logLevel)

	// Connect to database
	db, err := repository.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		l.Error("Failed to connect to database: %v", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			l.Error("Failed to close database: %v", err)
		}
	}()

	// Run migrations
	if err := repository.RunMigrations(db); err != nil {
		l.Error("Failed to run migrations: %v", err)
		os.Exit(1)
	}
	l.Info("Database migrations completed")

	// Create router
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	// CORS configuration
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3001"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(corsHandler.Handler)

	// Mock authentication middleware
	r.Use(middleware.MockAuthMiddleware)

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// API routes will be added here by the handlers
	r.Route("/api/v1", func(r chi.Router) {
		// Profile routes
		r.Route("/profile", func(r chi.Router) {
			// These will be implemented in Phase 3
			r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNotImplemented)
				_, _ = w.Write([]byte(`{"error": "Not implemented yet"}`))
			})
			r.Post("/generate", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNotImplemented)
				_, _ = w.Write([]byte(`{"error": "Not implemented yet"}`))
			})
			r.Post("/tailor", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNotImplemented)
				_, _ = w.Write([]byte(`{"error": "Not implemented yet"}`))
			})
			r.Get("/{profileId}/cv", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNotImplemented)
				_, _ = w.Write([]byte(`{"error": "Not implemented yet"}`))
			})
		})

		// Reference letter routes
		r.Route("/reference-letters", func(r chi.Router) {
			r.Post("/", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNotImplemented)
				_, _ = w.Write([]byte(`{"error": "Not implemented yet"}`))
			})
		})
	})

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	l.Info("Starting server on %s", addr)

	// Graceful shutdown
	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Error("Server error: %v", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	l.Info("Shutting down server...")
}
