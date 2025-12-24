package security

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrPasswordTooShort is returned when password is too short
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
)

// PasswordHandler handles password hashing and verification
type PasswordHandler struct {
	cost int
}

// NewPasswordHandler creates a new password handler
func NewPasswordHandler() *PasswordHandler {
	return &PasswordHandler{
		cost: bcrypt.DefaultCost,
	}
}

// Hash generates a bcrypt hash from a password
func (p *PasswordHandler) Hash(password string) (string, error) {
	if len(password) < 8 {
		return "", ErrPasswordTooShort
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), p.cost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// Verify compares a password with a hash
func (p *PasswordHandler) Verify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidatePassword validates password strength
func (p *PasswordHandler) ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}
