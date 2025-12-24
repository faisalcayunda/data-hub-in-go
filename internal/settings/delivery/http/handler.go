package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	settingsDomain "portal-data-backend/internal/settings/domain"
	"portal-data-backend/internal/settings/usecase"
	"portal-data-backend/infrastructure/http/response"
	pkgErrors "portal-data-backend/pkg/errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	settingsUsecase usecase.Usecase
	validator        *validator.Validate
}

func NewHandler(settingsUsecase usecase.Usecase) *Handler {
	return &Handler{
		settingsUsecase: settingsUsecase,
		validator:        validator.New(),
	}
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "Setting ID is required", nil)
		return
	}

	setting, err := h.settingsUsecase.GetByID(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Setting retrieved successfully", setting)
}

func (h *Handler) GetByKey(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	if key == "" {
		response.BadRequest(w, response.CodeBadRequest, "Setting key is required", nil)
		return
	}

	var userID *string
	if uid := r.URL.Query().Get("user_id"); uid != "" {
		userID = &uid
	}

	setting, err := h.settingsUsecase.GetByKey(r.Context(), key, userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Setting retrieved successfully", setting)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	req := &settingsDomain.ListSettingsRequest{
		Page:  parseIntQuery(r, "page", 1),
		Limit: parseIntQuery(r, "limit", 20),
		Search: r.URL.Query().Get("search"),
	}

	// Parse optional filters
	if category := r.URL.Query().Get("category"); category != "" {
		req.Category = &category
	}
	if userID := r.URL.Query().Get("user_id"); userID != "" {
		req.UserID = &userID
	}
	if settingType := r.URL.Query().Get("type"); settingType != "" {
		req.Type = &settingType
	}

	resp, err := h.settingsUsecase.List(r.Context(), req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Settings retrieved successfully", resp)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req settingsDomain.CreateSettingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, response.CodeBadRequest, "Invalid request body", nil)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		response.ValidationError(w, response.CodeValidationFailed, "Validation failed", h.formatValidationErrors(err))
		return
	}

	setting, err := h.settingsUsecase.Create(r.Context(), &req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.Created(w, response.CodeCreated, "Setting created successfully", setting)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "Setting ID is required", nil)
		return
	}

	var req settingsDomain.UpdateSettingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, response.CodeBadRequest, "Invalid request body", nil)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		response.ValidationError(w, response.CodeValidationFailed, "Validation failed", h.formatValidationErrors(err))
		return
	}

	setting, err := h.settingsUsecase.Update(r.Context(), id, &req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Setting updated successfully", setting)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "Setting ID is required", nil)
		return
	}

	if err := h.settingsUsecase.Delete(r.Context(), id); err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Setting deleted successfully", nil)
}

func (h *Handler) GetByKeys(w http.ResponseWriter, r *http.Request) {
	keysParam := r.URL.Query().Get("keys")
	if keysParam == "" {
		response.BadRequest(w, response.CodeBadRequest, "Keys parameter is required", nil)
		return
	}

	keys := strings.Split(keysParam, ",")
	for i := range keys {
		keys[i] = strings.TrimSpace(keys[i])
	}

	var userID *string
	if uid := r.URL.Query().Get("user_id"); uid != "" {
		userID = &uid
	}

	settings, err := h.settingsUsecase.GetByKeys(r.Context(), keys, userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Settings retrieved successfully", map[string]interface{}{"settings": settings})
}

func (h *Handler) GetByCategory(w http.ResponseWriter, r *http.Request) {
	category := chi.URLParam(r, "category")
	if category == "" {
		response.BadRequest(w, response.CodeBadRequest, "Category is required", nil)
		return
	}

	page := parseIntQuery(r, "page", 1)
	limit := parseIntQuery(r, "limit", 20)

	var userID *string
	if uid := r.URL.Query().Get("user_id"); uid != "" {
		userID = &uid
	}

	resp, err := h.settingsUsecase.GetByCategory(r.Context(), category, userID, page, limit)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Settings retrieved successfully", resp)
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	switch {
	case errors.Is(err, pkgErrors.ErrNotFound):
		response.NotFound(w, response.CodeNotFound, "Setting not found", nil)
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
	r.Route("/settings", func(r chi.Router) {
		r.Get("/", handler.List)
		r.Post("/", handler.Create)
		r.Get("/keys", handler.GetByKeys)
		r.Get("/category/{category}", handler.GetByCategory)
		r.Get("/key/{key}", handler.GetByKey)
		r.Get("/{id}", handler.GetByID)
		r.Put("/{id}", handler.Update)
		r.Delete("/{id}", handler.Delete)
	})
}
