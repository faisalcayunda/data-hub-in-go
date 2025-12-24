package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	orgDomain "portal-data-backend/internal/organization/domain"
	"portal-data-backend/internal/organization/usecase"
	"portal-data-backend/infrastructure/http/response"
	pkgErrors "portal-data-backend/pkg/errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// Handler handles HTTP requests for organization
type Handler struct {
	orgUsecase usecase.Usecase
	validator  *validator.Validate
}

// NewHandler creates a new organization handler
func NewHandler(orgUsecase usecase.Usecase) *Handler {
	return &Handler{
		orgUsecase: orgUsecase,
		validator:  validator.New(),
	}
}

// GetByID handles getting an organization by ID
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "Organization ID is required", nil)
		return
	}

	org, err := h.orgUsecase.GetByID(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Organization retrieved successfully", org)
}

// GetByCode handles getting an organization by code
func (h *Handler) GetByCode(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		response.BadRequest(w, response.CodeBadRequest, "Organization code is required", nil)
		return
	}

	org, err := h.orgUsecase.GetByCode(r.Context(), code)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Organization retrieved successfully", org)
}

// List handles listing organizations
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	req := &orgDomain.ListOrganizationsRequest{
		Page:      parseIntQuery(r, "page", 1),
		Limit:     parseIntQuery(r, "limit", 20),
		Status:    r.URL.Query().Get("status"),
		Search:    r.URL.Query().Get("search"),
		SortBy:    r.URL.Query().Get("sort_by"),
		SortOrder: r.URL.Query().Get("sort_order"),
	}

	resp, err := h.orgUsecase.List(r.Context(), req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Organizations retrieved successfully", resp)
}

// Create handles creating a new organization
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req orgDomain.CreateOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, response.CodeBadRequest, "Invalid request body", nil)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		response.ValidationError(w, response.CodeValidationFailed, "Validation failed", h.formatValidationErrors(err))
		return
	}

	// Get creator ID from context
	creatorID, _ := r.Context().Value("user_id").(string)

	org, err := h.orgUsecase.Create(r.Context(), &req, creatorID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.Created(w, response.CodeCreated, "Organization created successfully", org)
}

// Update handles updating an organization
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "Organization ID is required", nil)
		return
	}

	var req orgDomain.UpdateOrganizationRequest
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

	org, err := h.orgUsecase.Update(r.Context(), id, &req, updaterID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Organization updated successfully", org)
}

// Delete handles deleting an organization
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "Organization ID is required", nil)
		return
	}

	if err := h.orgUsecase.Delete(r.Context(), id); err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Organization deleted successfully", nil)
}

// UpdateStatus handles updating organization status
func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "Organization ID is required", nil)
		return
	}

	var req struct {
		Status orgDomain.OrgStatus `json:"status" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, response.CodeBadRequest, "Invalid request body", nil)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		response.ValidationError(w, response.CodeValidationFailed, "Validation failed", h.formatValidationErrors(err))
		return
	}

	if err := h.orgUsecase.UpdateStatus(r.Context(), id, req.Status); err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Organization status updated successfully", nil)
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	switch {
	case errors.Is(err, pkgErrors.ErrNotFound):
		response.NotFound(w, response.CodeNotFound, "Organization not found", nil)
	case errors.Is(err, pkgErrors.ErrAlreadyExists):
		response.Conflict(w, response.CodeConflict, "Organization code already exists", nil)
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

// RegisterRoutes registers organization routes
func RegisterRoutes(r chi.Router, handler *Handler) {
	r.Route("/organizations", func(r chi.Router) {
		r.Get("/", handler.List)
		r.Post("/", handler.Create)
		r.Get("/code/{code}", handler.GetByCode)
		r.Get("/{id}", handler.GetByID)
		r.Put("/{id}", handler.Update)
		r.Delete("/{id}", handler.Delete)
		r.Patch("/{id}/status", handler.UpdateStatus)
	})
}
