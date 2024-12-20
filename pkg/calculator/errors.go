package calculator

import (
	"fmt"
)

type ErrorType int // Custom error types

const (
	ErrInvalidCharacter ErrorType = iota
	ErrMismatchedParentheses
	ErrInsufficientValues
	ErrDivisionByZero
	ErrTooManyValues
	ErrTooLargeNumber
	ErrMismatchOperator
)

type CalcError struct {
	Type    ErrorType
	Message string
}

func (e CalcError) Error() string {
	return e.Message
}

func NewCalcError(errType ErrorType, details string) error {
	var message string

	switch errType {
	case ErrInvalidCharacter:
		message = fmt.Sprintf("invalid character: %s", details)
	case ErrMismatchedParentheses:
		message = fmt.Sprintf("mismatched parentheses: %s", details)
	case ErrInsufficientValues:
		message = fmt.Sprintf("insufficient values in expression: %s", details)
	case ErrDivisionByZero:
		message = "division by zero"
	case ErrTooManyValues:
		message = fmt.Sprintf("too many values in expression: %s", details)
	case ErrTooLargeNumber:
		message = fmt.Sprintf("number too large: %s", details)
	case ErrMismatchOperator:
		message = fmt.Sprintf("mismatched operator: %s", details)
	default:
		message = "unknown error"
	}

	return CalcError{Type: errType, Message: message}
}
