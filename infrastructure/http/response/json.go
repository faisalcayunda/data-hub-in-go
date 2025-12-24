package response

import (
	"encoding/json"
	"net/http"
)

// Response is the standard API response structure
type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta contains pagination metadata
type Meta struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}

// ErrorDetail contains error details
type ErrorDetail struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

// ErrorResponse is the standard error response structure
type ErrorResponse struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Details []ErrorDetail `json:"details,omitempty"`
}

// Response codes
const (
	CodeSuccess              = "OPERATION_SUCCESSFUL"
	CodeCreated              = "RESOURCE_CREATED"
	CodeUpdated              = "RESOURCE_UPDATED"
	CodeDeleted              = "RESOURCE_DELETED"
	CodeBadRequest           = "BAD_REQUEST"
	CodeUnauthorized         = "UNAUTHORIZED"
	CodeForbidden            = "FORBIDDEN"
	CodeNotFound             = "NOT_FOUND"
	CodeConflict             = "CONFLICT"
	CodeValidationFailed     = "VALIDATION_FAILED"
	CodeInternalServerError   = "INTERNAL_SERVER_ERROR"
	CodeServiceUnavailable   = "SERVICE_UNAVAILABLE"
	CodeTooManyRequests      = "TOO_MANY_REQUESTS"
)

// JSON sends a JSON response
func JSON(w http.ResponseWriter, statusCode int, code, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := Response{
		Code:    code,
		Message: message,
		Data:    data,
	}

	_ = json.NewEncoder(w).Encode(resp)
}

// JSONWithMeta sends a JSON response with pagination metadata
func JSONWithMeta(w http.ResponseWriter, statusCode int, code, message string, data interface{}, meta Meta) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := Response{
		Code:    code,
		Message: message,
		Data:    data,
		Meta:    &meta,
	}

	_ = json.NewEncoder(w).Encode(resp)
}

// Created sends a 201 Created response
func Created(w http.ResponseWriter, code, message string, data interface{}) {
	JSON(w, http.StatusCreated, code, message, data)
}

// OK sends a 200 OK response
func OK(w http.ResponseWriter, code, message string, data interface{}) {
	JSON(w, http.StatusOK, code, message, data)
}

// NoContent sends a 204 No Content response
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// BadRequest sends a 400 Bad Request response
func BadRequest(w http.ResponseWriter, code, message string, details []ErrorDetail) {
	Error(w, http.StatusBadRequest, code, message, details)
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(w http.ResponseWriter, code, message string, details []ErrorDetail) {
	Error(w, http.StatusUnauthorized, code, message, details)
}

// Forbidden sends a 403 Forbidden response
func Forbidden(w http.ResponseWriter, code, message string, details []ErrorDetail) {
	Error(w, http.StatusForbidden, code, message, details)
}

// NotFound sends a 404 Not Found response
func NotFound(w http.ResponseWriter, code, message string, details []ErrorDetail) {
	Error(w, http.StatusNotFound, code, message, details)
}

// Conflict sends a 409 Conflict response
func Conflict(w http.ResponseWriter, code, message string, details []ErrorDetail) {
	Error(w, http.StatusConflict, code, message, details)
}

// ValidationError sends a 422 Unprocessable Entity response
func ValidationError(w http.ResponseWriter, code, message string, details []ErrorDetail) {
	Error(w, http.StatusUnprocessableEntity, code, message, details)
}

// InternalError sends a 500 Internal Server Error response
func InternalError(w http.ResponseWriter, code, message string, details []ErrorDetail) {
	Error(w, http.StatusInternalServerError, code, message, details)
}

// Error sends an error response
func Error(w http.ResponseWriter, statusCode int, code, message string, details []ErrorDetail) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := ErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
	}

	_ = json.NewEncoder(w).Encode(resp)
}
