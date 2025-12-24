package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	pubDomain "portal-data-backend/internal/publication/domain"
	"portal-data-backend/internal/publication/usecase"
	"portal-data-backend/infrastructure/http/response"
	pkgErrors "portal-data-backend/pkg/errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	pubUsecase usecase.Usecase
	validator  *validator.Validate
}

func NewHandler(pubUsecase usecase.Usecase) *Handler {
	return &Handler{
		pubUsecase: pubUsecase,
		validator:  validator.New(),
	}
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "Publication ID is required", nil)
		return
	}

	pub, err := h.pubUsecase.GetByID(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Publication retrieved successfully", pub)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	req := &pubDomain.ListPublicationsRequest{
		Page:  parseIntQuery(r, "page", 1),
		Limit: parseIntQuery(r, "limit", 20),
		Search: r.URL.Query().Get("search"),
	}

	// Parse optional filters
	if datasetID := r.URL.Query().Get("dataset_id"); datasetID != "" {
		req.DatasetID = &datasetID
	}
	if organizationID := r.URL.Query().Get("organization_id"); organizationID != "" {
		req.OrganizationID = &organizationID
	}
	if status := r.URL.Query().Get("status"); status != "" {
		req.Status = &status
	}
	if isFeatured := r.URL.Query().Get("is_featured"); isFeatured != "" {
		featured := isFeatured == "true"
		req.IsFeatured = &featured
	}

	resp, err := h.pubUsecase.List(r.Context(), req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Publications retrieved successfully", resp)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req pubDomain.CreatePublicationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, response.CodeBadRequest, "Invalid request body", nil)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		response.ValidationError(w, response.CodeValidationFailed, "Validation failed", h.formatValidationErrors(err))
		return
	}

	userID, _ := r.Context().Value("user_id").(string)

	pub, err := h.pubUsecase.Create(r.Context(), &req, userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.Created(w, response.CodeCreated, "Publication created successfully", pub)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "Publication ID is required", nil)
		return
	}

	var req pubDomain.UpdatePublicationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, response.CodeBadRequest, "Invalid request body", nil)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		response.ValidationError(w, response.CodeValidationFailed, "Validation failed", h.formatValidationErrors(err))
		return
	}

	userID, _ := r.Context().Value("user_id").(string)

	pub, err := h.pubUsecase.Update(r.Context(), id, &req, userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Publication updated successfully", pub)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "Publication ID is required", nil)
		return
	}

	if err := h.pubUsecase.Delete(r.Context(), id); err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Publication deleted successfully", nil)
}

func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "Publication ID is required", nil)
		return
	}

	var req struct {
		Status string `json:"status" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, response.CodeBadRequest, "Invalid request body", nil)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		response.ValidationError(w, response.CodeValidationFailed, "Validation failed", h.formatValidationErrors(err))
		return
	}

	if err := h.pubUsecase.UpdateStatus(r.Context(), id, req.Status); err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Publication status updated successfully", nil)
}

func (h *Handler) IncrementDownloadCount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "Publication ID is required", nil)
		return
	}

	if err := h.pubUsecase.IncrementDownloadCount(r.Context(), id); err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Download count incremented successfully", nil)
}

func (h *Handler) GetByDatasetID(w http.ResponseWriter, r *http.Request) {
	datasetID := chi.URLParam(r, "datasetId")
	if datasetID == "" {
		response.BadRequest(w, response.CodeBadRequest, "Dataset ID is required", nil)
		return
	}

	page := parseIntQuery(r, "page", 1)
	limit := parseIntQuery(r, "limit", 20)

	resp, err := h.pubUsecase.GetByDatasetID(r.Context(), datasetID, page, limit)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Dataset publications retrieved successfully", resp)
}

func (h *Handler) GetByOrganizationID(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")
	if orgID == "" {
		response.BadRequest(w, response.CodeBadRequest, "Organization ID is required", nil)
		return
	}

	page := parseIntQuery(r, "page", 1)
	limit := parseIntQuery(r, "limit", 20)

	resp, err := h.pubUsecase.GetByOrganizationID(r.Context(), orgID, page, limit)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Organization publications retrieved successfully", resp)
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	switch {
	case errors.Is(err, pkgErrors.ErrNotFound):
		response.NotFound(w, response.CodeNotFound, "Publication not found", nil)
	default:
		response.InternalError(w, response.CodeInternalServerError, "Internal server error", nil)
	}
}

func (h *Handler) formatValidationErrors(err error) []response.ErrorDetail {
	var details []response.ErrorDetail
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			details = append(details, response.ErrorDetail{
				Field:   fieldErr.Field(),
				Message: h.getValidationErrorMessage(fieldErr),
			})
		}
	}
	return details
}

func (h *Handler) getValidationErrorMessage(fieldErr validator.FieldError) string {
	switch fieldErr.Tag() {
	case "required":
		return fieldErr.Field() + " is required"
	case "min":
		return fieldErr.Field() + " must be at least " + fieldErr.Param() + " characters"
	case "max":
		return fieldErr.Field() + " must be at most " + fieldErr.Param() + " characters"
	default:
		return fieldErr.Field() + " is invalid"
	}
}

func parseIntQuery(r *http.Request, key string, defaultValue int) int {
	if value := r.URL.Query().Get(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func RegisterRoutes(r chi.Router, handler *Handler) {
	r.Route("/publications", func(r chi.Router) {
		r.Get("/", handler.List)
		r.Post("/", handler.Create)
		r.Get("/dataset/{datasetId}", handler.GetByDatasetID)
		r.Get("/organization/{orgId}", handler.GetByOrganizationID)
		r.Get("/{id}", handler.GetByID)
		r.Put("/{id}", handler.Update)
		r.Delete("/{id}", handler.Delete)
		r.Patch("/{id}/status", handler.UpdateStatus)
		r.Post("/{id}/download", handler.IncrementDownloadCount)
	})
}
