package controller

import (
	"context"
)

type controller struct{}

type Controller interface {
	Calculate(ctx context.Context, expression string) (float64, error)
}

func New() Controller {
	return &controller{}
}
