package errors

import (
	"errors"
	"fmt"
)

// Sentinel errors untuk berbagai use case
var (
	// General errors
	ErrNotFound       = errors.New("resource not found")
	ErrAlreadyExists  = errors.New("resource already exists")
	ErrInvalidInput   = errors.New("invalid input")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrInternal       = errors.New("internal server error")
	ErrDatabase       = errors.New("database error")
	ErrValidation     = errors.New("validation error")

	// Auth specific errors
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
	ErrTokenRevoked       = errors.New("token revoked")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserDisabled       = errors.New("user is disabled")
	ErrEmailTaken         = errors.New("email already taken")
	ErrUsernameTaken      = errors.New("username already taken")

	// User specific errors
	ErrUserInactive = errors.New("user is inactive")

	// Organization specific errors
	ErrOrgNotFound     = errors.New("organization not found")
	ErrOrgInactive     = errors.New("organization is inactive")
	ErrInvalidOrgCode  = errors.New("invalid organization code")

	// Dataset specific errors
	ErrDatasetNotFound      = errors.New("dataset not found")
	ErrDatasetAccessDenied  = errors.New("access to dataset denied")
	ErrInvalidDatasetStatus = errors.New("invalid dataset status")
)

// Wrap wraps an error with context
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// Wrapf wraps an error with formatted context
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), err)
}

// Is checks if err is target
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in err's chain that matches target
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// New returns a new error
func New(message string) error {
	return errors.New(message)
}
