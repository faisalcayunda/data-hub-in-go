package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	datasetDomain "portal-data-backend/internal/dataset/domain"
	"portal-data-backend/internal/dataset/usecase"
	"portal-data-backend/infrastructure/http/response"
	pkgErrors "portal-data-backend/pkg/errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// Handler handles HTTP requests for dataset
type Handler struct {
	datasetUsecase usecase.Usecase
	validator      *validator.Validate
}

// NewHandler creates a new dataset handler
func NewHandler(datasetUsecase usecase.Usecase) *Handler {
	return &Handler{
		datasetUsecase: datasetUsecase,
		validator:      validator.New(),
	}
}

// GetByID handles getting a dataset by ID
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "Dataset ID is required", nil)
		return
	}

	dataset, err := h.datasetUsecase.GetByID(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Dataset retrieved successfully", dataset)
}

// GetBySlug handles getting a dataset by slug
func (h *Handler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.BadRequest(w, response.CodeBadRequest, "Dataset slug is required", nil)
		return
	}

	dataset, err := h.datasetUsecase.GetBySlug(r.Context(), slug)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Dataset retrieved successfully", dataset)
}

// List handles listing datasets
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	req := &datasetDomain.ListDatasetsRequest{
		Page:             parseIntQuery(r, "page", 1),
		Limit:            parseIntQuery(r, "limit", 20),
		OrganizationID:   r.URL.Query().Get("organization_id"),
		TopicID:          r.URL.Query().Get("topic_id"),
		BusinessFieldID:  r.URL.Query().Get("business_field_id"),
		TagID:            r.URL.Query().Get("tag_id"),
		Status:           r.URL.Query().Get("status"),
		ValidationStatus: r.URL.Query().Get("validation_status"),
		Classification:   r.URL.Query().Get("classification"),
		Search:           r.URL.Query().Get("search"),
		SortBy:           r.URL.Query().Get("sort_by"),
		SortOrder:        r.URL.Query().Get("sort_order"),
	}

	resp, err := h.datasetUsecase.List(r.Context(), req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Datasets retrieved successfully", resp)
}

// Create handles creating a new dataset
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req datasetDomain.CreateDatasetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, response.CodeBadRequest, "Invalid request body", nil)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		response.ValidationError(w, response.CodeValidationFailed, "Validation failed", h.formatValidationErrors(err))
		return
	}

	// Get creator ID and organization ID from context
	creatorID, _ := r.Context().Value("user_id").(string)
	orgID, _ := r.Context().Value("organization_id").(string)
	if orgID == "" {
		response.BadRequest(w, response.CodeBadRequest, "Organization ID is required", nil)
		return
	}

	dataset, err := h.datasetUsecase.Create(r.Context(), &req, creatorID, orgID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.Created(w, response.CodeCreated, "Dataset created successfully", dataset)
}

// Update handles updating a dataset
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "Dataset ID is required", nil)
		return
	}

	var req datasetDomain.UpdateDatasetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, response.CodeBadRequest, "Invalid request body", nil)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		response.ValidationError(w, response.CodeValidationFailed, "Validation failed", h.formatValidationErrors(err))
		return
	}

	// Get updater ID from context
	updaterID, _ := r.Context().Value("user_id").(string)

	dataset, err := h.datasetUsecase.Update(r.Context(), id, &req, updaterID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Dataset updated successfully", dataset)
}

// Delete handles deleting a dataset
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "Dataset ID is required", nil)
		return
	}

	if err := h.datasetUsecase.Delete(r.Context(), id); err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Dataset deleted successfully", nil)
}

// UpdateStatus handles updating dataset status
func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "Dataset ID is required", nil)
		return
	}

	var req struct {
		Status datasetDomain.DatasetStatus `json:"status" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, response.CodeBadRequest, "Invalid request body", nil)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		response.ValidationError(w, response.CodeValidationFailed, "Validation failed", h.formatValidationErrors(err))
		return
	}

	if err := h.datasetUsecase.UpdateStatus(r.Context(), id, req.Status); err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Dataset status updated successfully", nil)
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	switch {
	case errors.Is(err, pkgErrors.ErrNotFound):
		response.NotFound(w, response.CodeNotFound, "Dataset not found", nil)
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

// RegisterRoutes registers dataset routes
func RegisterRoutes(r chi.Router, handler *Handler) {
	r.Route("/datasets", func(r chi.Router) {
		r.Get("/", handler.List)
		r.Post("/", handler.Create)
		r.Get("/slug/{slug}", handler.GetBySlug)
		r.Get("/{id}", handler.GetByID)
		r.Put("/{id}", handler.Update)
		r.Delete("/{id}", handler.Delete)
		r.Patch("/{id}/status", handler.UpdateStatus)
	})
}
