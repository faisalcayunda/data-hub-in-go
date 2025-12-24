package http

import (
	"encoding/json"
	"net/http"

	"portal-data-backend/infrastructure/http/response"
	"portal-data-backend/internal/auth/usecase"
	"portal-data-backend/pkg/errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// Handler handles HTTP requests for auth
type Handler struct {
	authUsecase usecase.Usecase
	validator   *validator.Validate
}

// NewHandler creates a new auth handler
func NewHandler(authUsecase usecase.Usecase) *Handler {
	return &Handler{
		authUsecase: authUsecase,
		validator:   validator.New(),
	}
}

// Login handles user login
// @Summary Login
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /auth/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, response.CodeBadRequest, "Invalid request body", nil)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		response.ValidationError(w, response.CodeValidationFailed, "Validation failed", h.formatValidationErrors(err))
		return
	}

	authResp, err := h.authUsecase.Login(r.Context(), req.ToDomain())
	if err != nil {
		h.handleError(w, err)
		return
	}

	httpResp := &AuthResponse{}
	httpResp.FromDomain(authResp)

	response.OK(w, response.CodeSuccess, "Login successful", httpResp)
}

// Register handles user registration
// @Summary Register
// @Description Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration data"
// @Success 201 {object} AuthResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Router /auth/register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, response.CodeBadRequest, "Invalid request body", nil)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		response.ValidationError(w, response.CodeValidationFailed, "Validation failed", h.formatValidationErrors(err))
		return
	}

	authResp, err := h.authUsecase.Register(r.Context(), req.ToDomain())
	if err != nil {
		h.handleError(w, err)
		return
	}

	httpResp := &AuthResponse{}
	httpResp.FromDomain(authResp)

	response.Created(w, response.CodeCreated, "Registration successful", httpResp)
}

// Logout handles user logout
// @Summary Logout
// @Description Logout user and revoke tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LogoutRequest true "Refresh token"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /auth/logout [post]
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var req LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, response.CodeBadRequest, "Invalid request body", nil)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		response.ValidationError(w, response.CodeValidationFailed, "Validation failed", h.formatValidationErrors(err))
		return
	}

	// Get access token from Authorization header
	accessToken := r.Header.Get("Authorization")
	if len(accessToken) > 7 && accessToken[:7] == "Bearer " {
		accessToken = accessToken[7:]
	}

	if err := h.authUsecase.Logout(r.Context(), accessToken, req.RefreshToken); err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Logout successful", MessageResponse{Message: "Successfully logged out"})
}

// RefreshToken handles token refresh
// @Summary Refresh Token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /auth/refresh [post]
func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, response.CodeBadRequest, "Invalid request body", nil)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		response.ValidationError(w, response.CodeValidationFailed, "Validation failed", h.formatValidationErrors(err))
		return
	}

	authResp, err := h.authUsecase.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		h.handleError(w, err)
		return
	}

	httpResp := &AuthResponse{}
	httpResp.FromDomain(authResp)

	response.OK(w, response.CodeSuccess, "Token refreshed successfully", httpResp)
}

// RevokeAllTokens handles revoking all user tokens
// @Summary Revoke All Tokens
// @Description Revoke all tokens for the current user
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} MessageResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /auth/revoke-all [post]
func (h *Handler) RevokeAllTokens(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID := r.Context().Value("user_id").(string)
	if userID == "" {
		response.Unauthorized(w, response.CodeUnauthorized, "Unauthorized", nil)
		return
	}

	if err := h.authUsecase.RevokeAllTokens(r.Context(), userID); err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "All tokens revoked successfully", MessageResponse{Message: "All tokens revoked"})
}

// GetCurrentUser handles getting current user
// @Summary Get Current User
// @Description Get current user information
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} UserInfo
// @Failure 401 {object} response.ErrorResponse
// @Router /me [get]
func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID := r.Context().Value("user_id").(string)
	if userID == "" {
		response.Unauthorized(w, response.CodeUnauthorized, "Unauthorized", nil)
		return
	}

	userInfo, err := h.authUsecase.GetCurrentUser(r.Context(), userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	httpResp := UserInfo{
		ID:             userInfo.ID,
		OrganizationID: userInfo.OrganizationID,
		RoleID:         userInfo.RoleID,
		Name:           userInfo.Name,
		Username:       userInfo.Username,
		Email:          userInfo.Email,
		Thumbnail:      userInfo.Thumbnail,
	}

	response.OK(w, response.CodeSuccess, "User retrieved successfully", httpResp)
}

// handleError handles errors and returns appropriate HTTP responses
func (h *Handler) handleError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	switch {
	case errors.Is(err, errors.ErrInvalidCredentials):
		response.Unauthorized(w, response.CodeUnauthorized, "Invalid credentials", nil)
	case errors.Is(err, errors.ErrUserDisabled):
		response.Forbidden(w, response.CodeForbidden, "User account is disabled", nil)
	case errors.Is(err, errors.ErrEmailTaken):
		response.Conflict(w, response.CodeConflict, "Email already registered", nil)
	case errors.Is(err, errors.ErrUsernameTaken):
		response.Conflict(w, response.CodeConflict, "Username already taken", nil)
	case errors.Is(err, errors.ErrInvalidToken):
		response.Unauthorized(w, response.CodeUnauthorized, "Invalid token", nil)
	case errors.Is(err, errors.ErrTokenExpired):
		response.Unauthorized(w, response.CodeUnauthorized, "Token expired", nil)
	case errors.Is(err, errors.ErrTokenRevoked):
		response.Unauthorized(w, response.CodeUnauthorized, "Token revoked", nil)
	case errors.Is(err, errors.ErrNotFound):
		response.NotFound(w, response.CodeNotFound, "Resource not found", nil)
	default:
		response.InternalError(w, response.CodeInternalServerError, "Internal server error", nil)
	}
}

// formatValidationErrors formats validation errors into ErrorDetail slice
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

// getValidationErrorMessage returns a user-friendly validation error message
func (h *Handler) getValidationErrorMessage(fieldErr validator.FieldError) string {
	switch fieldErr.Tag() {
	case "required":
		return fieldErr.Field() + " is required"
	case "email":
		return fieldErr.Field() + " must be a valid email"
	case "min":
		return fieldErr.Field() + " must be at least " + fieldErr.Param() + " characters"
	case "max":
		return fieldErr.Field() + " must be at most " + fieldErr.Param() + " characters"
	case "alphanum":
		return fieldErr.Field() + " must contain only alphanumeric characters"
	default:
		return fieldErr.Field() + " is invalid"
	}
}

// RegisterRoutes registers auth routes
func RegisterRoutes(r chi.Router, handler *Handler) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", handler.Login)
		r.Post("/register", handler.Register)
		r.Post("/logout", handler.Logout)
		r.Post("/refresh", handler.RefreshToken)
		r.Post("/revoke-all", handler.RevokeAllTokens)
	})

	r.Get("/me", handler.GetCurrentUser)
}
