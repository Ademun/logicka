package ast

import "fmt"

type ValidationError struct {
	message string
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{message: message}
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s", e.message)
}
