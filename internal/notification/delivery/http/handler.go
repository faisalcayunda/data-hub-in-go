package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	notifDomain "portal-data-backend/internal/notification/domain"
	"portal-data-backend/internal/notification/usecase"
	"portal-data-backend/infrastructure/http/response"
	pkgErrors "portal-data-backend/pkg/errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	notifUsecase usecase.Usecase
	validator    *validator.Validate
}

func NewHandler(notifUsecase usecase.Usecase) *Handler {
	return &Handler{
		notifUsecase: notifUsecase,
		validator:    validator.New(),
	}
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "Notification ID is required", nil)
		return
	}

	notif, err := h.notifUsecase.GetByID(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Notification retrieved successfully", notif)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	req := &notifDomain.ListNotificationsRequest{
		Page:  parseIntQuery(r, "page", 1),
		Limit: parseIntQuery(r, "limit", 20),
	}

	// Get user ID from context
	userID, _ := r.Context().Value("user_id").(string)
	req.UserID = &userID

	// Parse optional filters
	if notifType := r.URL.Query().Get("type"); notifType != "" {
		req.Type = &notifType
	}
	if category := r.URL.Query().Get("category"); category != "" {
		req.Category = &category
	}
	if isRead := r.URL.Query().Get("is_read"); isRead != "" {
		read := isRead == "true"
		req.IsRead = &read
	}
	if startDate := r.URL.Query().Get("start_date"); startDate != "" {
		req.StartDate = &startDate
	}
	if endDate := r.URL.Query().Get("end_date"); endDate != "" {
		req.EndDate = &endDate
	}

	resp, err := h.notifUsecase.List(r.Context(), req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Notifications retrieved successfully", resp)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req notifDomain.CreateNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, response.CodeBadRequest, "Invalid request body", nil)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		response.ValidationError(w, response.CodeValidationFailed, "Validation failed", h.formatValidationErrors(err))
		return
	}

	notif, err := h.notifUsecase.Create(r.Context(), &req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.Created(w, response.CodeCreated, "Notification created successfully", notif)
}

func (h *Handler) BulkCreate(w http.ResponseWriter, r *http.Request) {
	var req notifDomain.BulkCreateNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, response.CodeBadRequest, "Invalid request body", nil)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		response.ValidationError(w, response.CodeValidationFailed, "Validation failed", h.formatValidationErrors(err))
		return
	}

	if err := h.notifUsecase.BulkCreate(r.Context(), &req); err != nil {
		h.handleError(w, err)
		return
	}

	response.Created(w, response.CodeCreated, "Notifications created successfully", nil)
}

func (h *Handler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	var req notifDomain.MarkAsReadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, response.CodeBadRequest, "Invalid request body", nil)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		response.ValidationError(w, response.CodeValidationFailed, "Validation failed", h.formatValidationErrors(err))
		return
	}

	userID, _ := r.Context().Value("user_id").(string)

	if err := h.notifUsecase.MarkAsRead(r.Context(), req.NotificationIDs, userID); err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Notifications marked as read successfully", nil)
}

func (h *Handler) MarkAllAsRead(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("user_id").(string)

	if err := h.notifUsecase.MarkAllAsRead(r.Context(), userID); err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "All notifications marked as read successfully", nil)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "Notification ID is required", nil)
		return
	}

	if err := h.notifUsecase.Delete(r.Context(), id); err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Notification deleted successfully", nil)
}

func (h *Handler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("user_id").(string)

	count, err := h.notifUsecase.GetUnreadCount(r.Context(), userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Unread count retrieved successfully", notifDomain.UnreadCountResponse{Count: count})
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	switch {
	case errors.Is(err, pkgErrors.ErrNotFound):
		response.NotFound(w, response.CodeNotFound, "Notification not found", nil)
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
	r.Route("/notifications", func(r chi.Router) {
		r.Get("/", handler.List)
		r.Post("/", handler.Create)
		r.Post("/bulk", handler.BulkCreate)
		r.Post("/mark-read", handler.MarkAsRead)
		r.Post("/mark-all-read", handler.MarkAllAsRead)
		r.Get("/unread-count", handler.GetUnreadCount)
		r.Get("/{id}", handler.GetByID)
		r.Delete("/{id}", handler.Delete)
	})
}
