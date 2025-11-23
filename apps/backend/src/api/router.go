package api

import (
	"net/http"

	"github.com/credfolio/apps/backend/src/api/handlers"
	"github.com/credfolio/apps/backend/src/api/middleware"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

func NewRouter(uploadHandler *handlers.UploadHandler, profileHandler *handlers.ProfileHandler, cvHandler *handlers.CVHandler, tailorHandler *handlers.TailorHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(middleware.RequestLogger)
	r.Use(middleware.ErrorHandler)
	r.Use(cors.AllowAll().Handler)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	r.Post("/api/upload", uploadHandler.HandleUpload)
	r.Get("/api/profile", profileHandler.GetProfile)
	r.Get("/api/profile/cv", cvHandler.DownloadCV)
	r.Post("/api/profile/tailor", tailorHandler.TailorProfile)

	return r
}
