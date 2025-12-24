package middleware

import (
	"context"
	"net/http"
	"strings"

	"portal-data-backend/infrastructure/http/response"
	"portal-data-backend/infrastructure/security"
	"portal-data-backend/pkg/errors"
)

// Auth middleware validates JWT tokens
func Auth(jwtManager *security.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Unauthorized(w, response.CodeUnauthorized, "Authorization header required", nil)
				return
			}

			// Check Bearer prefix
			if !strings.HasPrefix(authHeader, "Bearer ") {
				response.Unauthorized(w, response.CodeUnauthorized, "Invalid authorization header format", nil)
				return
			}

			// Extract token
			token := authHeader[7:]

			// Validate token
			claims, err := jwtManager.ValidateToken(token)
			if err != nil {
				if errors.Is(err, errors.ErrTokenExpired) {
					response.Unauthorized(w, response.CodeUnauthorized, "Token expired", nil)
					return
				}
				response.Unauthorized(w, response.CodeUnauthorized, "Invalid token", nil)
				return
			}

			// Add user info to context
			ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
			ctx = context.WithValue(ctx, "organization_id", claims.OrganizationID)
			ctx = context.WithValue(ctx, "role_id", claims.RoleID)
			ctx = context.WithValue(ctx, "email", claims.Email)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
