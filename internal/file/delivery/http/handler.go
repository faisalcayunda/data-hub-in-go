package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	fileDomain "portal-data-backend/internal/file/domain"
	"portal-data-backend/internal/file/usecase"
	"portal-data-backend/infrastructure/http/response"
	pkgErrors "portal-data-backend/pkg/errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// readCloser wraps a bytes.Reader to implement multipart.File interface
type readCloser struct {
	*bytes.Reader
}

func (rc *readCloser) Close() error {
	return nil
}

type Handler struct {
	fileUsecase usecase.Usecase
	validator   *validator.Validate
}

func NewHandler(fileUsecase usecase.Usecase) *Handler {
	return &Handler{
		fileUsecase: fileUsecase,
		validator:   validator.New(),
	}
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "File ID is required", nil)
		return
	}

	file, err := h.fileUsecase.GetByID(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "File retrieved successfully", file)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	req := &fileDomain.ListFilesRequest{
		Page:  parseIntQuery(r, "page", 1),
		Limit: parseIntQuery(r, "limit", 20),
		Search: r.URL.Query().Get("search"),
	}

	// Parse optional filters
	if datasetID := r.URL.Query().Get("dataset_id"); datasetID != "" {
		req.DatasetID = &datasetID
	}
	if status := r.URL.Query().Get("status"); status != "" {
		req.Status = &status
	}

	resp, err := h.fileUsecase.List(r.Context(), req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Files retrieved successfully", resp)
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (max 32MB)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		response.BadRequest(w, response.CodeBadRequest, "Failed to parse form data", nil)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		response.BadRequest(w, response.CodeBadRequest, "File is required", nil)
		return
	}
	defer file.Close()

	// Get form fields
	fileName := r.FormValue("filename")
	if fileName == "" {
		fileName = header.Filename
	}

	mimeType := r.FormValue("mime_type")
	if mimeType == "" {
		mimeType = header.Header.Get("Content-Type")
	}

	fileSize := header.Size
	if fileSize == 0 {
		// Read file to get size if not available
		var buffer bytes.Buffer
		fileSize, err = io.Copy(&buffer, file)
		if err != nil {
			response.BadRequest(w, response.CodeBadRequest, "Failed to read file", nil)
			return
		}
		// Wrap buffer with a type that implements multipart.File
		file = &readCloser{Reader: bytes.NewReader(buffer.Bytes())}
	}

	var datasetID *string
	if dsID := r.FormValue("dataset_id"); dsID != "" {
		datasetID = &dsID
	}

	// Get user ID from context
	userID, _ := r.Context().Value("user_id").(string)

	uploadResp, err := h.fileUsecase.Upload(r.Context(), fileName, fileSize, mimeType, file, datasetID, userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.Created(w, response.CodeCreated, "File uploaded successfully", uploadResp)
}

func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "File ID is required", nil)
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

	if err := h.fileUsecase.UpdateStatus(r.Context(), id, fileDomain.FileStatus(req.Status)); err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "File status updated successfully", nil)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, response.CodeBadRequest, "File ID is required", nil)
		return
	}

	if err := h.fileUsecase.Delete(r.Context(), id); err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "File deleted successfully", nil)
}

func (h *Handler) GetByDatasetID(w http.ResponseWriter, r *http.Request) {
	datasetID := chi.URLParam(r, "datasetId")
	if datasetID == "" {
		response.BadRequest(w, response.CodeBadRequest, "Dataset ID is required", nil)
		return
	}

	page := parseIntQuery(r, "page", 1)
	limit := parseIntQuery(r, "limit", 20)

	resp, err := h.fileUsecase.GetByDatasetID(r.Context(), datasetID, page, limit)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Dataset files retrieved successfully", resp)
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	switch {
	case errors.Is(err, pkgErrors.ErrNotFound):
		response.NotFound(w, response.CodeNotFound, "File not found", nil)
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
		return fmt.Sprintf("%s is required", fieldErr.Field())
	default:
		return fmt.Sprintf("%s is invalid", fieldErr.Field())
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
	r.Route("/files", func(r chi.Router) {
		r.Get("/", handler.List)
		r.Post("/upload", handler.Upload)
		r.Get("/{id}", handler.GetByID)
		r.Patch("/{id}/status", handler.UpdateStatus)
		r.Delete("/{id}", handler.Delete)
		r.Get("/dataset/{datasetId}", handler.GetByDatasetID)
	})
}
