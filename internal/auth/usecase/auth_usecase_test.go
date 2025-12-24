package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"portal-data-backend/infrastructure/config"
	"portal-data-backend/infrastructure/security"
	"portal-data-backend/internal/auth/domain"
	"portal-data-backend/internal/auth/usecase"
	pkgerrors "portal-data-backend/pkg/errors"

	"github.com/google/uuid"
)

// mockUserRepository is a mock implementation of UserRepository
type mockUserRepository struct {
	users map[string]*domain.User

	getUserByIDFunc       func(ctx context.Context, id string) (*domain.User, error)
	getUserByEmailFunc    func(ctx context.Context, email string) (*domain.User, error)
	getUserByUsernameFunc func(ctx context.Context, username string) (*domain.User, error)
	createUserFunc        func(ctx context.Context, user *domain.User) error
	isEmailExistsFunc     func(ctx context.Context, email string) (bool, error)
	isUsernameExistsFunc  func(ctx context.Context, username string) (bool, error)
}

func (m *mockUserRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	if m.getUserByIDFunc != nil {
		return m.getUserByIDFunc(ctx, id)
	}
	user, ok := m.users[id]
	if !ok {
		return nil, pkgerrors.ErrNotFound
	}
	return user, nil
}

func (m *mockUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.getUserByEmailFunc != nil {
		return m.getUserByEmailFunc(ctx, email)
	}
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, pkgerrors.ErrNotFound
}

func (m *mockUserRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	if m.getUserByUsernameFunc != nil {
		return m.getUserByUsernameFunc(ctx, username)
	}
	for _, user := range m.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, pkgerrors.ErrNotFound
}

func (m *mockUserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	if m.createUserFunc != nil {
		return m.createUserFunc(ctx, user)
	}
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	return nil
}

func (m *mockUserRepository) DeleteUser(ctx context.Context, id string) error {
	return nil
}

func (m *mockUserRepository) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, int, error) {
	return nil, 0, nil
}

func (m *mockUserRepository) IsEmailExists(ctx context.Context, email string) (bool, error) {
	if m.isEmailExistsFunc != nil {
		return m.isEmailExistsFunc(ctx, email)
	}
	for _, user := range m.users {
		if user.Email == email {
			return true, nil
		}
	}
	return false, nil
}

func (m *mockUserRepository) IsUsernameExists(ctx context.Context, username string) (bool, error) {
	if m.isUsernameExistsFunc != nil {
		return m.isUsernameExistsFunc(ctx, username)
	}
	for _, user := range m.users {
		if user.Username == username {
			return true, nil
		}
	}
	return false, nil
}

// mockTokenRepository is a mock implementation of TokenRepository
type mockTokenRepository struct {
	tokens map[string]*domain.Token

	createTokenFunc       func(ctx context.Context, token *domain.Token) error
	getTokenByRefreshFunc func(ctx context.Context, refreshToken string) (*domain.Token, error)
	getTokenByAccessFunc  func(ctx context.Context, accessToken string) (*domain.Token, error)
	revokeTokenFunc       func(ctx context.Context, id string) error
	revokeUserTokensFunc  func(ctx context.Context, userID string) error
}

func (m *mockTokenRepository) CreateToken(ctx context.Context, token *domain.Token) error {
	if m.createTokenFunc != nil {
		return m.createTokenFunc(ctx, token)
	}
	m.tokens[token.ID] = token
	return nil
}

func (m *mockTokenRepository) GetTokenByRefreshToken(ctx context.Context, refreshToken string) (*domain.Token, error) {
	if m.getTokenByRefreshFunc != nil {
		return m.getTokenByRefreshFunc(ctx, refreshToken)
	}
	for _, token := range m.tokens {
		if token.RefreshToken == refreshToken {
			return token, nil
		}
	}
	return nil, pkgerrors.ErrNotFound
}

func (m *mockTokenRepository) GetTokenByAccessToken(ctx context.Context, accessToken string) (*domain.Token, error) {
	if m.getTokenByAccessFunc != nil {
		return m.getTokenByAccessFunc(ctx, accessToken)
	}
	for _, token := range m.tokens {
		if token.AccessToken == accessToken {
			return token, nil
		}
	}
	return nil, pkgerrors.ErrNotFound
}

func (m *mockTokenRepository) RevokeToken(ctx context.Context, id string) error {
	if m.revokeTokenFunc != nil {
		return m.revokeTokenFunc(ctx, id)
	}
	if token, ok := m.tokens[id]; ok {
		token.Revoked = true
		return nil
	}
	return pkgerrors.ErrNotFound
}

func (m *mockTokenRepository) RevokeUserTokens(ctx context.Context, userID string) error {
	if m.revokeUserTokensFunc != nil {
		return m.revokeUserTokensFunc(ctx, userID)
	}
	for _, token := range m.tokens {
		if token.UserID == userID {
			token.Revoked = true
		}
	}
	return nil
}

func (m *mockTokenRepository) DeleteToken(ctx context.Context, id string) error {
	return nil
}

func (m *mockTokenRepository) CleanupExpiredTokens(ctx context.Context) error {
	return nil
}

// Helper function to create a test user with hashed password
func createTestUser(id, email, password string) (*domain.User, error) {
	hasher := security.NewPasswordHandler()
	hash, err := hasher.Hash(password)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:             id,
		OrganizationID: uuid.New().String(),
		RoleID:         uuid.New().String(),
		Name:           "Test User",
		Username:       "testuser",
		Email:          email,
		PasswordHash:   hash,
		Status:         domain.UserStatusActive,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}

// Test Login Success
func TestLogin_Success(t *testing.T) {
	ctx := context.Background()

	// Create test user
	user, err := createTestUser(uuid.New().String(), "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Setup mocks
	userRepo := &mockUserRepository{
		users: map[string]*domain.User{user.ID: user},
		getUserByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
			if email == user.Email {
				return user, nil
			}
			return nil, pkgerrors.ErrNotFound
		},
	}

	tokenRepo := &mockTokenRepository{
		tokens: make(map[string]*domain.Token),
	}

	jwtConfig := &config.JWTConfig{
		Secret:             "test-secret-key-for-testing",
		AccessTokenExpiry:  time.Hour,
		RefreshTokenExpiry: 24 * time.Hour,
		Issuer:             "test",
	}
	jwtManager := security.NewJWTManager(jwtConfig)
	passwordHasher := security.NewPasswordHandler()

	// Create usecase
	authUsecase := usecase.NewAuthUsecase(userRepo, tokenRepo, jwtManager, passwordHasher)

	// Execute
	req := &domain.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	resp, err := authUsecase.Login(ctx, req)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	if resp.User.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, resp.User.Email)
	}

	if resp.AccessToken == "" {
		t.Error("Expected access token, got empty string")
	}

	if resp.RefreshToken == "" {
		t.Error("Expected refresh token, got empty string")
	}

	if resp.TokenType != "Bearer" {
		t.Errorf("Expected token type Bearer, got %s", resp.TokenType)
	}
}

// Test Login Invalid Credentials
func TestLogin_InvalidCredentials(t *testing.T) {
	ctx := context.Background()

	// Create test user
	user, err := createTestUser(uuid.New().String(), "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Setup mocks
	userRepo := &mockUserRepository{
		users: map[string]*domain.User{user.ID: user},
		getUserByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
			return user, nil
		},
	}

	tokenRepo := &mockTokenRepository{
		tokens: make(map[string]*domain.Token),
	}

	jwtConfig := &config.JWTConfig{
		Secret:             "test-secret-key-for-testing",
		AccessTokenExpiry:  time.Hour,
		RefreshTokenExpiry: 24 * time.Hour,
		Issuer:             "test",
	}
	jwtManager := security.NewJWTManager(jwtConfig)
	passwordHasher := security.NewPasswordHandler()

	// Create usecase
	authUsecase := usecase.NewAuthUsecase(userRepo, tokenRepo, jwtManager, passwordHasher)

	// Execute with wrong password
	req := &domain.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	resp, err := authUsecase.Login(ctx, req)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if !errors.Is(err, pkgerrors.ErrInvalidCredentials) {
		t.Errorf("Expected ErrInvalidCredentials, got %v", err)
	}

	if resp != nil {
		t.Error("Expected nil response, got response")
	}
}

// Test Register Success
func TestRegister_Success(t *testing.T) {
	ctx := context.Background()

	// Setup mocks
	userRepo := &mockUserRepository{
		users: make(map[string]*domain.User),
		isEmailExistsFunc: func(ctx context.Context, email string) (bool, error) {
			return false, nil
		},
		isUsernameExistsFunc: func(ctx context.Context, username string) (bool, error) {
			return false, nil
		},
		createUserFunc: func(ctx context.Context, user *domain.User) error {
			return nil
		},
	}

	tokenRepo := &mockTokenRepository{
		tokens: make(map[string]*domain.Token),
	}

	jwtConfig := &config.JWTConfig{
		Secret:             "test-secret-key-for-testing",
		AccessTokenExpiry:  time.Hour,
		RefreshTokenExpiry: 24 * time.Hour,
		Issuer:             "test",
	}
	jwtManager := security.NewJWTManager(jwtConfig)
	passwordHasher := security.NewPasswordHandler()

	// Create usecase
	authUsecase := usecase.NewAuthUsecase(userRepo, tokenRepo, jwtManager, passwordHasher)

	// Execute
	req := &domain.RegisterRequest{
		OrganizationID: uuid.New().String(),
		RoleID:         uuid.New().String(),
		Name:           "New User",
		Username:       "newuser",
		Email:          "new@example.com",
		Password:       "password123",
	}

	resp, err := authUsecase.Register(ctx, req)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	if resp.User.Email != req.Email {
		t.Errorf("Expected email %s, got %s", req.Email, resp.User.Email)
	}

	if resp.AccessToken == "" {
		t.Error("Expected access token, got empty string")
	}
}

// Test Register Email Already Exists
func TestRegister_EmailAlreadyExists(t *testing.T) {
	ctx := context.Background()

	// Setup mocks
	userRepo := &mockUserRepository{
		users: make(map[string]*domain.User),
		isEmailExistsFunc: func(ctx context.Context, email string) (bool, error) {
			return true, nil
		},
	}

	tokenRepo := &mockTokenRepository{
		tokens: make(map[string]*domain.Token),
	}

	jwtConfig := &config.JWTConfig{
		Secret:             "test-secret-key-for-testing",
		AccessTokenExpiry:  time.Hour,
		RefreshTokenExpiry: 24 * time.Hour,
		Issuer:             "test",
	}
	jwtManager := security.NewJWTManager(jwtConfig)
	passwordHasher := security.NewPasswordHandler()

	// Create usecase
	authUsecase := usecase.NewAuthUsecase(userRepo, tokenRepo, jwtManager, passwordHasher)

	// Execute
	req := &domain.RegisterRequest{
		OrganizationID: uuid.New().String(),
		RoleID:         uuid.New().String(),
		Name:           "New User",
		Username:       "newuser",
		Email:          "existing@example.com",
		Password:       "password123",
	}

	resp, err := authUsecase.Register(ctx, req)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if !errors.Is(err, pkgerrors.ErrEmailTaken) {
		t.Errorf("Expected ErrEmailTaken, got %v", err)
	}

	if resp != nil {
		t.Error("Expected nil response, got response")
	}
}

// Test RefreshToken Success
func TestRefreshToken_Success(t *testing.T) {
	ctx := context.Background()

	// Create test user
	user, err := createTestUser(uuid.New().String(), "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Setup mocks
	userRepo := &mockUserRepository{
		users: map[string]*domain.User{user.ID: user},
		getUserByIDFunc: func(ctx context.Context, id string) (*domain.User, error) {
			return user, nil
		},
	}

	jwtConfig := &config.JWTConfig{
		Secret:             "test-secret-key-for-testing",
		AccessTokenExpiry:  time.Hour,
		RefreshTokenExpiry: 24 * time.Hour,
		Issuer:             "test",
	}
	jwtManager := security.NewJWTManager(jwtConfig)

	// Generate initial token pair
	tokenPair, err := jwtManager.GenerateTokenPair(user.ID, user.OrganizationID, user.RoleID, user.Email)
	if err != nil {
		t.Fatalf("Failed to generate token pair: %v", err)
	}

	// Create stored token
	storedToken := &domain.Token{
		ID:           uuid.New().String(),
		UserID:       user.ID,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		Revoked:      false,
		CreatedAt:    time.Now(),
	}

	tokenRepo := &mockTokenRepository{
		tokens: map[string]*domain.Token{storedToken.ID: storedToken},
		getTokenByRefreshFunc: func(ctx context.Context, refreshToken string) (*domain.Token, error) {
			if refreshToken == tokenPair.RefreshToken {
				return storedToken, nil
			}
			return nil, pkgerrors.ErrNotFound
		},
		revokeTokenFunc: func(ctx context.Context, id string) error {
			return nil
		},
		createTokenFunc: func(ctx context.Context, token *domain.Token) error {
			return nil
		},
	}

	passwordHasher := security.NewPasswordHandler()

	// Create usecase
	authUsecase := usecase.NewAuthUsecase(userRepo, tokenRepo, jwtManager, passwordHasher)

	// Execute
	resp, err := authUsecase.RefreshToken(ctx, tokenPair.RefreshToken)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	if resp.AccessToken == "" {
		t.Error("Expected new access token, got empty string")
	}
}
