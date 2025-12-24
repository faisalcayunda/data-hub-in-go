package usecase

import (
	"context"
	"fmt"
	"time"

	"portal-data-backend/internal/auth/domain"
	"portal-data-backend/infrastructure/security"
	"portal-data-backend/pkg/errors"

	"github.com/google/uuid"
)

// authUsecase implements the Usecase interface
type authUsecase struct {
	userRepo       domain.UserRepository
	tokenRepo      domain.TokenRepository
	jwtManager     *security.JWTManager
	passwordHasher *security.PasswordHandler
}

// NewAuthUsecase creates a new auth usecase
func NewAuthUsecase(
	userRepo domain.UserRepository,
	tokenRepo domain.TokenRepository,
	jwtManager *security.JWTManager,
	passwordHasher *security.PasswordHandler,
) Usecase {
	return &authUsecase{
		userRepo:       userRepo,
		tokenRepo:      tokenRepo,
		jwtManager:     jwtManager,
		passwordHasher: passwordHasher,
	}
}

// Login authenticates a user and returns tokens
func (a *authUsecase) Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error) {
	// Get user by email
	user, err := a.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, errors.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Verify password
	if !a.passwordHasher.Verify(req.Password, user.PasswordHash) {
		return nil, errors.ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive() {
		return nil, errors.ErrUserDisabled
	}

	// Generate tokens
	tokenPair, err := a.jwtManager.GenerateTokenPair(
		user.ID,
		user.OrganizationID,
		user.RoleID,
		user.Email,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Store refresh token in database
	token := &domain.Token{
		ID:           uuid.New().String(),
		UserID:       user.ID,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    time.Now().Add(24 * time.Hour * 7), // 7 days
		Revoked:      false,
		CreatedAt:    time.Now(),
	}

	if err := a.tokenRepo.CreateToken(ctx, token); err != nil {
		return nil, fmt.Errorf("failed to store token: %w", err)
	}

	return &domain.AuthResponse{
		User:         user.ToUserInfo(),
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		TokenType:    tokenPair.TokenType,
	}, nil
}

// Register creates a new user account
func (a *authUsecase) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	// Check if email already exists
	exists, err := a.userRepo.IsEmailExists(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, errors.ErrEmailTaken
	}

	// Check if username already exists
	exists, err = a.userRepo.IsUsernameExists(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username existence: %w", err)
	}
	if exists {
		return nil, errors.ErrUsernameTaken
	}

	// Hash password
	passwordHash, err := a.passwordHasher.Hash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &domain.User{
		ID:             uuid.New().String(),
		OrganizationID: req.OrganizationID,
		RoleID:         req.RoleID,
		Name:           req.Name,
		Username:       req.Username,
		Email:          req.Email,
		PasswordHash:   passwordHash,
		Status:         domain.UserStatusActive,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if req.EmployeeID != "" {
		user.EmployeeID = &req.EmployeeID
	}
	if req.Position != "" {
		user.Position = &req.Position
	}
	if req.Address != "" {
		user.Address = &req.Address
	}
	if req.Phone != "" {
		user.Phone = &req.Phone
	}

	if err := a.userRepo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens
	tokenPair, err := a.jwtManager.GenerateTokenPair(
		user.ID,
		user.OrganizationID,
		user.RoleID,
		user.Email,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Store refresh token in database
	token := &domain.Token{
		ID:           uuid.New().String(),
		UserID:       user.ID,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    time.Now().Add(24 * time.Hour * 7),
		Revoked:      false,
		CreatedAt:    time.Now(),
	}

	if err := a.tokenRepo.CreateToken(ctx, token); err != nil {
		return nil, fmt.Errorf("failed to store token: %w", err)
	}

	return &domain.AuthResponse{
		User:         user.ToUserInfo(),
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		TokenType:    tokenPair.TokenType,
	}, nil
}

// Logout logs out a user by revoking their tokens
func (a *authUsecase) Logout(ctx context.Context, accessToken, refreshToken string) error {
	// Get token by refresh token
	token, err := a.tokenRepo.GetTokenByRefreshToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return errors.ErrInvalidToken
		}
		return fmt.Errorf("failed to get token: %w", err)
	}

	// Verify access token matches
	if token.AccessToken != accessToken {
		return errors.ErrInvalidToken
	}

	// Revoke token
	if err := a.tokenRepo.RevokeToken(ctx, token.ID); err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}

	return nil
}

// RefreshToken refreshes an access token using a refresh token
func (a *authUsecase) RefreshToken(ctx context.Context, refreshToken string) (*domain.AuthResponse, error) {
	// Get stored token
	storedToken, err := a.tokenRepo.GetTokenByRefreshToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, errors.ErrInvalidToken
		}
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	// Check if token is valid
	if !storedToken.IsValid() {
		return nil, errors.ErrTokenExpired
	}

	// Get user
	user, err := a.userRepo.GetUserByID(ctx, storedToken.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Generate new token pair
	tokenPair, err := a.jwtManager.GenerateTokenPair(
		user.ID,
		user.OrganizationID,
		user.RoleID,
		user.Email,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new tokens: %w", err)
	}

	// Revoke old token
	if err := a.tokenRepo.RevokeToken(ctx, storedToken.ID); err != nil {
		return nil, fmt.Errorf("failed to revoke old token: %w", err)
	}

	// Store new refresh token
	newToken := &domain.Token{
		ID:           uuid.New().String(),
		UserID:       user.ID,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    time.Now().Add(24 * time.Hour * 7),
		Revoked:      false,
		CreatedAt:    time.Now(),
	}

	if err := a.tokenRepo.CreateToken(ctx, newToken); err != nil {
		return nil, fmt.Errorf("failed to store new token: %w", err)
	}

	return &domain.AuthResponse{
		User:         user.ToUserInfo(),
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		TokenType:    tokenPair.TokenType,
	}, nil
}

// RevokeAllTokens revokes all tokens for a user
func (a *authUsecase) RevokeAllTokens(ctx context.Context, userID string) error {
	if err := a.tokenRepo.RevokeUserTokens(ctx, userID); err != nil {
		return fmt.Errorf("failed to revoke user tokens: %w", err)
	}
	return nil
}

// ValidateToken validates a token and returns the claims
func (a *authUsecase) ValidateToken(ctx context.Context, token string) (*domain.TokenClaims, error) {
	claims, err := a.jwtManager.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	// Check if token is stored and not revoked
	storedToken, err := a.tokenRepo.GetTokenByAccessToken(ctx, token)
	if err != nil && !errors.Is(err, errors.ErrNotFound) {
		return nil, fmt.Errorf("failed to get stored token: %w", err)
	}

	if storedToken != nil && !storedToken.IsValid() {
		return nil, errors.ErrTokenRevoked
	}

	return &domain.TokenClaims{
		UserID:         claims.UserID,
		OrganizationID: claims.OrganizationID,
		RoleID:         claims.RoleID,
		Email:          claims.Email,
	}, nil
}

// GetCurrentUser retrieves the current user by ID
func (a *authUsecase) GetCurrentUser(ctx context.Context, userID string) (*domain.UserInfo, error) {
	user, err := a.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	info := user.ToUserInfo()
	return &info, nil
}
