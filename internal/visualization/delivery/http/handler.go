package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	vizDomain "portal-data-backend/internal/visualization/domain"
	"portal-data-backend/internal/visualization/usecase"
	"portal-data-backend/infrastructure/http/response"
	pkgErrors "portal-data-backend/pkg/errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	vizUsecase usecase.Usecase
	validator  *validator.Validate
}

func NewHandler(vizUsecase usecase.Usecase) *Handler {
	return &Handler{
		vizUsecase: vizUsecase,
		validator:  validator.New(),
	}
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "Visualization ID is required", nil)
		return
	}

	viz, err := h.vizUsecase.GetByID(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Visualization retrieved successfully", viz)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	req := &vizDomain.ListVisualizationsRequest{
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
	if topicID := r.URL.Query().Get("topic_id"); topicID != "" {
		req.TopicID = &topicID
	}
	if vizType := r.URL.Query().Get("type"); vizType != "" {
		req.Type = &vizType
	}
	if status := r.URL.Query().Get("status"); status != "" {
		req.Status = &status
	}
	if isHighlight := r.URL.Query().Get("is_highlight"); isHighlight != "" {
		highlight := isHighlight == "true"
		req.IsHighlight = &highlight
	}

	resp, err := h.vizUsecase.List(r.Context(), req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Visualizations retrieved successfully", resp)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req vizDomain.CreateVisualizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, response.CodeBadRequest, "Invalid request body", nil)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		response.ValidationError(w, response.CodeValidationFailed, "Validation failed", h.formatValidationErrors(err))
		return
	}

	userID, _ := r.Context().Value("user_id").(string)

	viz, err := h.vizUsecase.Create(r.Context(), &req, userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.Created(w, response.CodeCreated, "Visualization created successfully", viz)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "Visualization ID is required", nil)
		return
	}

	var req vizDomain.UpdateVisualizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, response.CodeBadRequest, "Invalid request body", nil)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		response.ValidationError(w, response.CodeValidationFailed, "Validation failed", h.formatValidationErrors(err))
		return
	}

	userID, _ := r.Context().Value("user_id").(string)

	viz, err := h.vizUsecase.Update(r.Context(), id, &req, userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Visualization updated successfully", viz)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "Visualization ID is required", nil)
		return
	}

	if err := h.vizUsecase.Delete(r.Context(), id); err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Visualization deleted successfully", nil)
}

func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "Visualization ID is required", nil)
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

	if err := h.vizUsecase.UpdateStatus(r.Context(), id, req.Status); err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Visualization status updated successfully", nil)
}

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.vizUsecase.GetStats(r.Context())
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Visualization stats retrieved successfully", stats)
}

func (h *Handler) GetByDatasetID(w http.ResponseWriter, r *http.Request) {
	datasetID := chi.URLParam(r, "datasetId")
	if datasetID == "" {
		response.BadRequest(w, response.CodeBadRequest, "Dataset ID is required", nil)
		return
	}

	page := parseIntQuery(r, "page", 1)
	limit := parseIntQuery(r, "limit", 20)

	resp, err := h.vizUsecase.GetByDatasetID(r.Context(), datasetID, page, limit)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Dataset visualizations retrieved successfully", resp)
}

func (h *Handler) GetByOrganizationID(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgId")
	if orgID == "" {
		response.BadRequest(w, response.CodeBadRequest, "Organization ID is required", nil)
		return
	}

	page := parseIntQuery(r, "page", 1)
	limit := parseIntQuery(r, "limit", 20)

	resp, err := h.vizUsecase.GetByOrganizationID(r.Context(), orgID, page, limit)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Organization visualizations retrieved successfully", resp)
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	switch {
	case errors.Is(err, pkgErrors.ErrNotFound):
		response.NotFound(w, response.CodeNotFound, "Visualization not found", nil)
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
	r.Route("/visualizations", func(r chi.Router) {
		r.Get("/", handler.List)
		r.Post("/", handler.Create)
		r.Get("/stats", handler.GetStats)
		r.Get("/dataset/{datasetId}", handler.GetByDatasetID)
		r.Get("/organization/{orgId}", handler.GetByOrganizationID)
		r.Get("/{id}", handler.GetByID)
		r.Put("/{id}", handler.Update)
		r.Delete("/{id}", handler.Delete)
		r.Patch("/{id}/status", handler.UpdateStatus)
	})
}
