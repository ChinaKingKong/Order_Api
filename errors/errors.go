package errors

import (
	"errors"
	"fmt"
)

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrInvalidOrderStatus = errors.New("invalid order status")
	ErrInvalidOrderID     = errors.New("invalid order ID")
	ErrEmptyOrder        = errors.New("order is empty")
	ErrInvalidUserID     = errors.New("invalid user ID")
	ErrInvalidAmount     = errors.New("invalid amount")
	ErrInvalidQuantity   = errors.New("invalid quantity")
	ErrInvalidPrice      = errors.New("invalid price")
	ErrDatabaseError     = errors.New("database error")
	ErrCacheError        = errors.New("cache error")
	ErrUnauthorized      = errors.New("unauthorized access")
	ErrForbidden         = errors.New("forbidden")
)

type AppError struct {
	Err     error
	Message string
	Code    int
	Stack   string
}

func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(message string) error {
	return &AppError{
		Message: message,
		Code:    500,
	}
}

func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return &AppError{
		Err:     err,
		Message: fmt.Sprintf("%s: %v", message, err),
		Code:    500,
	}
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}