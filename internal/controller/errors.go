package controller

import "fmt"

type ErrorType string

const (
	ErrRequest ErrorType = "request error"
	ErrServer  ErrorType = "server error"
)

type CtrlError struct {
	Err  error
	Type ErrorType
}

func (c CtrlError) Error() string {
	return fmt.Sprintf("%s: %s", c.Type, c.Err.Error())
}

func NewRequestError(err error) CtrlError {
	return CtrlError{
		Err:  err,
		Type: ErrRequest,
	}
}

func NewServerError(err error) CtrlError {
	return CtrlError{
		Err:  err,
		Type: ErrServer,
	}
}
