package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	integrationDomain "portal-data-backend/internal/integration/domain"
	"portal-data-backend/internal/integration/usecase"
	"portal-data-backend/infrastructure/http/response"
	pkgErrors "portal-data-backend/pkg/errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	integrationUsecase usecase.Usecase
	validator           *validator.Validate
}

func NewHandler(integrationUsecase usecase.Usecase) *Handler {
	return &Handler{
		integrationUsecase: integrationUsecase,
		validator:           validator.New(),
	}
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "Integration ID is required", nil)
		return
	}

	integration, err := h.integrationUsecase.GetByID(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Integration retrieved successfully", integration)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	req := &integrationDomain.ListIntegrationsRequest{
		Page:  parseIntQuery(r, "page", 1),
		Limit: parseIntQuery(r, "limit", 20),
		Search: r.URL.Query().Get("search"),
	}

	// Parse optional filters
	if orgID := r.URL.Query().Get("organization_id"); orgID != "" {
		req.OrganizationID = &orgID
	}
	if intType := r.URL.Query().Get("type"); intType != "" {
		req.Type = &intType
	}
	if status := r.URL.Query().Get("status"); status != "" {
		req.Status = &status
	}

	resp, err := h.integrationUsecase.List(r.Context(), req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Integrations retrieved successfully", resp)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req integrationDomain.CreateIntegrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, response.CodeBadRequest, "Invalid request body", nil)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		response.ValidationError(w, response.CodeValidationFailed, "Validation failed", h.formatValidationErrors(err))
		return
	}

	userID, _ := r.Context().Value("user_id").(string)

	integration, err := h.integrationUsecase.Create(r.Context(), &req, userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.Created(w, response.CodeCreated, "Integration created successfully", integration)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "Integration ID is required", nil)
		return
	}

	var req integrationDomain.UpdateIntegrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, response.CodeBadRequest, "Invalid request body", nil)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		response.ValidationError(w, response.CodeValidationFailed, "Validation failed", h.formatValidationErrors(err))
		return
	}

	integration, err := h.integrationUsecase.Update(r.Context(), id, &req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Integration updated successfully", integration)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "Integration ID is required", nil)
		return
	}

	if err := h.integrationUsecase.Delete(r.Context(), id); err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Integration deleted successfully", nil)
}

func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "Integration ID is required", nil)
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

	if err := h.integrationUsecase.UpdateStatus(r.Context(), id, req.Status); err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Integration status updated successfully", nil)
}

func (h *Handler) Sync(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "Integration ID is required", nil)
		return
	}

	if err := h.integrationUsecase.Sync(r.Context(), id); err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Integration synced successfully", nil)
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	switch {
	case errors.Is(err, pkgErrors.ErrNotFound):
		response.NotFound(w, response.CodeNotFound, "Integration not found", nil)
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

func RegisterRoutes(r chi.Router, handler *Handler) {
	r.Route("/integrations", func(r chi.Router) {
		r.Get("/", handler.List)
		r.Post("/", handler.Create)
		r.Get("/{id}", handler.GetByID)
		r.Put("/{id}", handler.Update)
		r.Delete("/{id}", handler.Delete)
		r.Patch("/{id}/status", handler.UpdateStatus)
		r.Post("/{id}/sync", handler.Sync)
	})
}
