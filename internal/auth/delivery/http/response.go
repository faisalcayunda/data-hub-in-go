package http

import (
	"portal-data-backend/internal/auth/domain"
)

// AuthResponse represents HTTP response for authentication
type AuthResponse struct {
	User         UserInfo `json:"user"`
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int64    `json:"expires_in"`
	TokenType    string   `json:"token_type"`
}

// FromDomain converts domain response to HTTP response
func (r *AuthResponse) FromDomain(resp *domain.AuthResponse) {
	r.User = UserInfo{
		ID:             resp.User.ID,
		OrganizationID: resp.User.OrganizationID,
		RoleID:         resp.User.RoleID,
		Name:           resp.User.Name,
		Username:       resp.User.Username,
		Email:          resp.User.Email,
		Thumbnail:      resp.User.Thumbnail,
	}
	r.AccessToken = resp.AccessToken
	r.RefreshToken = resp.RefreshToken
	r.ExpiresIn = resp.ExpiresIn
	r.TokenType = resp.TokenType
}

// UserInfo represents user information in HTTP response
type UserInfo struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	RoleID         string `json:"role_id"`
	Name           string `json:"name"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	Thumbnail      string `json:"thumbnail,omitempty"`
}

// MessageResponse represents a simple message response
type MessageResponse struct {
	Message string `json:"message"`
}
