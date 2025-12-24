package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	dataRowDomain "portal-data-backend/internal/data_row/domain"
	"portal-data-backend/internal/data_row/usecase"
	"portal-data-backend/infrastructure/http/response"
	pkgErrors "portal-data-backend/pkg/errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	dataRowUsecase usecase.Usecase
	validator       *validator.Validate
}

func NewHandler(dataRowUsecase usecase.Usecase) *Handler {
	return &Handler{
		dataRowUsecase: dataRowUsecase,
		validator:       validator.New(),
	}
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "Data row ID is required", nil)
		return
	}

	row, err := h.dataRowUsecase.GetByID(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Data row retrieved successfully", row)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	datasetID := chi.URLParam(r, "datasetId")
	if datasetID == "" {
		response.BadRequest(w, response.CodeBadRequest, "Dataset ID is required", nil)
		return
	}

	req := &dataRowDomain.ListDataRowsRequest{
		Page:      parseIntQuery(r, "page", 1),
		Limit:     parseIntQuery(r, "limit", 20),
		DatasetID: datasetID,
		Search:    r.URL.Query().Get("search"),
	}

	resp, err := h.dataRowUsecase.List(r.Context(), req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Data rows retrieved successfully", resp)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req dataRowDomain.CreateDataRowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, response.CodeBadRequest, "Invalid request body", nil)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		response.ValidationError(w, response.CodeValidationFailed, "Validation failed", h.formatValidationErrors(err))
		return
	}

	userID, _ := r.Context().Value("user_id").(string)

	row, err := h.dataRowUsecase.Create(r.Context(), &req, userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.Created(w, response.CodeCreated, "Data row created successfully", row)
}

func (h *Handler) BulkCreate(w http.ResponseWriter, r *http.Request) {
	var req dataRowDomain.BulkCreateDataRowsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, response.CodeBadRequest, "Invalid request body", nil)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		response.ValidationError(w, response.CodeValidationFailed, "Validation failed", h.formatValidationErrors(err))
		return
	}

	userID, _ := r.Context().Value("user_id").(string)

	if err := h.dataRowUsecase.BulkCreate(r.Context(), &req, userID); err != nil {
		h.handleError(w, err)
		return
	}

	response.Created(w, response.CodeCreated, "Data rows created successfully", nil)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "Data row ID is required", nil)
		return
	}

	var req dataRowDomain.UpdateDataRowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, response.CodeBadRequest, "Invalid request body", nil)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		response.ValidationError(w, response.CodeValidationFailed, "Validation failed", h.formatValidationErrors(err))
		return
	}

	row, err := h.dataRowUsecase.Update(r.Context(), id, &req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Data row updated successfully", row)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "Data row ID is required", nil)
		return
	}

	if err := h.dataRowUsecase.Delete(r.Context(), id); err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Data row deleted successfully", nil)
}

func (h *Handler) DeleteByDatasetID(w http.ResponseWriter, r *http.Request) {
	datasetID := chi.URLParam(r, "datasetId")
	if datasetID == "" {
		response.BadRequest(w, response.CodeBadRequest, "Dataset ID is required", nil)
		return
	}

	if err := h.dataRowUsecase.DeleteByDatasetID(r.Context(), datasetID); err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Data rows deleted successfully", nil)
}

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	datasetID := chi.URLParam(r, "datasetId")
	if datasetID == "" {
		response.BadRequest(w, response.CodeBadRequest, "Dataset ID is required", nil)
		return
	}

	stats, err := h.dataRowUsecase.GetStats(r.Context(), datasetID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Data row stats retrieved successfully", stats)
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	switch {
	case errors.Is(err, pkgErrors.ErrNotFound):
		response.NotFound(w, response.CodeNotFound, "Data row not found", nil)
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
	r.Route("/datasets/{datasetId}/data-rows", func(r chi.Router) {
		r.Get("/", handler.List)
		r.Post("/", handler.Create)
		r.Post("/bulk", handler.BulkCreate)
		r.Get("/stats", handler.GetStats)
		r.Delete("/", handler.DeleteByDatasetID)
	})
	r.Route("/data-rows", func(r chi.Router) {
		r.Get("/{id}", handler.GetByID)
		r.Put("/{id}", handler.Update)
		r.Delete("/{id}", handler.Delete)
	})
}
