package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	userDomain "portal-data-backend/internal/user/domain"
	"portal-data-backend/internal/user/usecase"
	"portal-data-backend/infrastructure/http/response"
	pkgErrors "portal-data-backend/pkg/errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// Handler handles HTTP requests for user
type Handler struct {
	userUsecase usecase.Usecase
	validator   *validator.Validate
}

// NewHandler creates a new user handler
func NewHandler(userUsecase usecase.Usecase) *Handler {
	return &Handler{
		userUsecase: userUsecase,
		validator:   validator.New(),
	}
}

// GetUserByID handles getting a user by ID
func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		response.BadRequest(w, response.CodeBadRequest, "User ID is required", nil)
		return
	}

	userInfo, err := h.userUsecase.GetUserByID(r.Context(), userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "User retrieved successfully", userInfo)
}

// ListUsers handles listing users with pagination
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	req := &userDomain.ListUsersRequest{
		Page:      parseIntQuery(r, "page", 1),
		Limit:     parseIntQuery(r, "limit", 20),
		Search:    r.URL.Query().Get("search"),
		Status:    r.URL.Query().Get("status"),
		SortBy:    r.URL.Query().Get("sort_by"),
		SortOrder: r.URL.Query().Get("sort_order"),
	}

	resp, err := h.userUsecase.ListUsers(r.Context(), req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Users retrieved successfully", resp)
}

// UpdateUser handles updating a user
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		response.BadRequest(w, response.CodeBadRequest, "User ID is required", nil)
		return
	}

	var req userDomain.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, response.CodeBadRequest, "Invalid request body", nil)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		response.ValidationError(w, response.CodeValidationFailed, "Validation failed", h.formatValidationErrors(err))
		return
	}

	userInfo, err := h.userUsecase.UpdateUser(r.Context(), userID, &req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "User updated successfully", userInfo)
}

// DeleteUser handles deleting a user
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		response.BadRequest(w, response.CodeBadRequest, "User ID is required", nil)
		return
	}

	if err := h.userUsecase.DeleteUser(r.Context(), userID); err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "User deleted successfully", map[string]string{"id": userID})
}

// UpdateUserStatus handles updating user status
func (h *Handler) UpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		response.BadRequest(w, response.CodeBadRequest, "User ID is required", nil)
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

	if err := h.userUsecase.UpdateUserStatus(r.Context(), userID, req.Status); err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "User status updated successfully", nil)
}

// handleError handles errors
func (h *Handler) handleError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	switch {
	case errors.Is(err, pkgErrors.ErrNotFound):
		response.NotFound(w, response.CodeNotFound, "User not found", nil)
	default:
		response.InternalError(w, response.CodeInternalServerError, "Internal server error", nil)
	}
}

// formatValidationErrors formats validation errors
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

// getValidationErrorMessage returns validation error message
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

// parseIntQuery parses integer query parameter with default value
func parseIntQuery(r *http.Request, key string, defaultValue int) int {
	if value := r.URL.Query().Get(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// RegisterRoutes registers user routes
func RegisterRoutes(r chi.Router, handler *Handler) {
	r.Route("/users", func(r chi.Router) {
		r.Get("/", handler.ListUsers)
		r.Get("/{id}", handler.GetUserByID)
		r.Put("/{id}", handler.UpdateUser)
		r.Delete("/{id}", handler.DeleteUser)
		r.Patch("/{id}/status", handler.UpdateUserStatus)
	})
}
