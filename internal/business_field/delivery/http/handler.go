package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	bfDomain "portal-data-backend/internal/business_field/domain"
	"portal-data-backend/internal/business_field/usecase"
	"portal-data-backend/infrastructure/http/response"
	pkgErrors "portal-data-backend/pkg/errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	bfUsecase usecase.Usecase
	validator  *validator.Validate
}

func NewHandler(bfUsecase usecase.Usecase) *Handler {
	return &Handler{
		bfUsecase: bfUsecase,
		validator: validator.New(),
	}
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "Business Field ID is required", nil)
		return
	}

	bf, err := h.bfUsecase.GetByID(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Business field retrieved successfully", bf)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	req := &bfDomain.ListBusinessFieldsRequest{
		Page:  parseIntQuery(r, "page", 1),
		Limit: parseIntQuery(r, "limit", 20),
		Search: r.URL.Query().Get("search"),
	}

	resp, err := h.bfUsecase.List(r.Context(), req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Business fields retrieved successfully", resp)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req bfDomain.CreateBusinessFieldRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, response.CodeBadRequest, "Invalid request body", nil)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		response.ValidationError(w, response.CodeValidationFailed, "Validation failed", h.formatValidationErrors(err))
		return
	}

	bf, err := h.bfUsecase.Create(r.Context(), &req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.Created(w, response.CodeCreated, "Business field created successfully", bf)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "Business Field ID is required", nil)
		return
	}

	var req bfDomain.UpdateBusinessFieldRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, response.CodeBadRequest, "Invalid request body", nil)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		response.ValidationError(w, response.CodeValidationFailed, "Validation failed", h.formatValidationErrors(err))
		return
	}

	bf, err := h.bfUsecase.Update(r.Context(), id, &req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Business field updated successfully", bf)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "Business Field ID is required", nil)
		return
	}

	if err := h.bfUsecase.Delete(r.Context(), id); err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Business field deleted successfully", nil)
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	switch {
	case errors.Is(err, pkgErrors.ErrNotFound):
		response.NotFound(w, response.CodeNotFound, "Business field not found", nil)
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
	r.Route("/business-fields", func(r chi.Router) {
		r.Get("/", handler.List)
		r.Post("/", handler.Create)
		r.Get("/{id}", handler.GetByID)
		r.Put("/{id}", handler.Update)
		r.Delete("/{id}", handler.Delete)
	})
}
