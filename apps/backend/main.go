package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/credfolio/apps/backend/api/generated"
	"github.com/credfolio/apps/backend/internal/handler"
	authmiddleware "github.com/credfolio/apps/backend/internal/handler/middleware"
	"github.com/credfolio/apps/backend/internal/repository"
	"github.com/credfolio/apps/backend/internal/service"
	"github.com/credfolio/apps/backend/pkg/ai"
	"github.com/credfolio/apps/backend/pkg/config"
	"github.com/credfolio/apps/backend/pkg/logger"
	"github.com/credfolio/apps/backend/pkg/pdf"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	// Configure logger: timestamps with microseconds and short file info
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

	// Load environment from .env files in development if present
	_ = godotenv.Load()
	_ = godotenv.Load(".env.development")

	cfg := config.Load()
	port := cfg.Port
	log.Printf("starting backend (pid=%d) on port %s", os.Getpid(), port)

	// Initialize database
	if err := repository.InitDB(cfg); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	// Run migrations
	if err := repository.RunMigrations(); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	db := repository.GetDB()

	// Initialize repositories
	profileRepo := repository.NewGormProfileRepository(db)
	workExpRepo := repository.NewGormWorkExperienceRepository(db)
	credibilityRepo := repository.NewGormCredibilityHighlightRepository(db)
	referenceLetterRepo := repository.NewGormReferenceLetterRepository(db)

	// Initialize providers
	appLogger := logger.New()
	llmProvider, err := ai.NewOpenAIProvider(cfg.OpenAIKey, "")
	if err != nil {
		log.Fatalf("failed to initialize AI provider: %v", err)
	}
	pdfExtractor := pdf.NewExtractor()

	// Initialize repositories
	jobMatchRepo := repository.NewGormJobMatchRepository(repository.GetDB())

	// Initialize services
	profileService := service.NewProfileService(
		profileRepo,
		workExpRepo,
		credibilityRepo,
		referenceLetterRepo,
		llmProvider,
		pdfExtractor,
		appLogger,
	)

	tailoringService := service.NewTailoringService(
		profileRepo,
		jobMatchRepo,
		llmProvider,
		appLogger,
	)

	// Initialize API handler
	apiHandler := handler.NewAPI(
		profileService,
		tailoringService,
		referenceLetterRepo,
		pdfExtractor,
		appLogger,
	)

	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	// Custom concise request logger with latency and status code
	r.Use(requestLogger)
	r.Use(chimiddleware.Recoverer)
	r.Use(cors.AllowAll().Handler)
	r.Use(authmiddleware.MockAuth) // Mock authentication middleware

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// Register generated API routes
	generated.HandlerFromMux(apiHandler, r)

	s := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("backend listening on :%s", port)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

// statusRecorder wraps http.ResponseWriter to capture the status code for logging.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// requestLogger logs method, path, status and latency for each request.
func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		latency := time.Since(start)
		log.Printf("%s %s -> %d (%s) ip=%s ua=%q",
			r.Method, r.URL.Path, rec.status, latency, r.RemoteAddr, r.UserAgent())
	})
}
