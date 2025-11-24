package handler

import (
	"net/http"

	"github.com/credfolio/apps/backend/api/generated"
	"github.com/credfolio/apps/backend/internal/repository"
	"github.com/credfolio/apps/backend/internal/service"
	"github.com/credfolio/apps/backend/pkg/logger"
	"github.com/credfolio/apps/backend/pkg/pdf"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// API implements generated.ServerInterface
type API struct {
	ProfileService      *service.ProfileService
	TailoringService    *service.TailoringService
	ReferenceLetterRepo repository.ReferenceLetterRepository
	PDFExtractor        pdf.ExtractorInterface
	Logger              *logger.Logger
}

func NewAPI(
	profileService *service.ProfileService,
	tailoringService *service.TailoringService,
	referenceLetterRepo repository.ReferenceLetterRepository,
	pdfExtractor pdf.ExtractorInterface,
	logger *logger.Logger,
) *API {
	return &API{
		ProfileService:      profileService,
		TailoringService:    tailoringService,
		ReferenceLetterRepo: referenceLetterRepo,
		PDFExtractor:        pdfExtractor,
		Logger:              logger,
	}
}

// Ensure API implements generated.ServerInterface
var _ generated.ServerInterface = (*API)(nil)

// TailorProfile implements generated.ServerInterface
// Implementation is in profile_handler.go

// DownloadCV implements generated.ServerInterface
func (a *API) DownloadCV(w http.ResponseWriter, r *http.Request, profileId openapi_types.UUID) {
	w.WriteHeader(http.StatusNotImplemented)
}
